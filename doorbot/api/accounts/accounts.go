// Package accounts handles account-related api endpoints
package accounts

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/auth"
	"github.com/masom/doorbot/doorbot/security"
	"github.com/masom/doorbot/doorbot/services/notifications"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

// AccountsViewModel represent a list of accounts
type AccountsViewModel struct {
	Accounts []*doorbot.Account `json:"accounts"`
}

// AccountViewModel used by api responses and requests. Wraps the Account object.
type AccountViewModel struct {
	Account *doorbot.Account `json:"account" binding:"required"`
}

// RegisterViewModel wraps the AccountRegister model.
type RegisterViewModel struct {
	Account AccountRegisterRequest `json:"account" binding:"required"`
}

// PublicAccountViewModel represent a public account VM
type PublicAccountViewModel struct {
	Account PublicAccount `json:"account"`
}

// PublicAccount holds publicly visible data
type PublicAccount struct {
	ID   uint
	Name string
	Host string
}

// AccountRegisterRequest represents an account request.
type AccountRegisterRequest struct {
	Host               string `json:"host"`
	Name               string `json:"name" binding:"required"`
	ContactName        string `json:"contact_name" binding:"required"`
	ContactPhoneNumber string `json:"contact_phone_number" binding:"required"`
	ContactEmail       string `json:"contact_email" binding:"required"`
}

// Validate the RegisterModel instance
func (vm RegisterViewModel) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	vm.Account.Name = strings.TrimSpace(vm.Account.Name)
	vm.Account.ContactName = strings.TrimSpace(vm.Account.ContactName)
	vm.Account.ContactEmail = strings.TrimSpace(vm.Account.ContactEmail)
	vm.Account.ContactPhoneNumber = strings.TrimSpace(vm.Account.ContactPhoneNumber)
	vm.Account.Host = strings.TrimSpace(vm.Account.Host)

	if len(vm.Account.Host) > 20 {
		errors.Add([]string{"host"}, "", "The host must be smaller than 20 characters")
	}

	if len(vm.Account.Name) < 3 || len(vm.Account.Name) > 100 {
		errors.Add([]string{"name"}, "", "The account name must be at least 3 characters or smaller than 255")
	}

	if len(vm.Account.ContactName) < 3 {
		errors.Add([]string{"contact_name"}, "", "The contact name must be at least 3 characters or smaller than 255")
	}

	if len(vm.Account.ContactPhoneNumber) < 10 || len(vm.Account.ContactPhoneNumber) > 50 {
		errors.Add([]string{"contact_phone_number"}, "", "A valid phone number must be provided")
	}

	_, err := mail.ParseAddress(vm.Account.ContactEmail)
	if len(vm.Account.ContactEmail) < 3 || err != nil {
		errors.Add([]string{"email"}, "", "A valid email address must be provided")
	}

	return errors
}

// Index action
func Index(render render.Render, r doorbot.Repositories, administrator *doorbot.Administrator) {
	repo := r.AccountRepository()
	accounts, err := repo.All(r.DB())

	if err != nil {
		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	render.JSON(http.StatusOK, AccountsViewModel{Accounts: accounts})
}

// Get return a specific account
func Get(render render.Render, r doorbot.Repositories, params martini.Params, session *auth.Authorization) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.AccountRepository()

	account, err := repo.Find(r.DB(), uint(id))

	if err != nil {
		log.WithFields(log.Fields{
			"account_id": id,
			"error":      err,
		}).Error("Api::Accounts->Get database error.")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if account == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{}))
		return
	}

	// Switch the view model depending on who/what requests the information.
	switch session.Type {
	case auth.AuthorizationAdministrator:
		render.JSON(http.StatusOK, AccountViewModel{Account: account})
	case auth.AuthorizationPerson:
		if session.Person.IsAccountManager() {
			render.JSON(http.StatusOK, AccountViewModel{Account: account})
			return
		}

		// Display a reduced version of the account.
		public := PublicAccount{
			ID:   account.ID,
			Name: account.Name,
			Host: account.Host,
		}

		render.JSON(http.StatusOK, PublicAccountViewModel{Account: public})
	default:
		render.Status(http.StatusForbidden)
		return
	}
}

// Register a new account ( used by the dashboard )
func Register(render render.Render, config *doorbot.DoorbotConfig, r doorbot.Repositories, n notifications.Notificator, data RegisterViewModel) {
	repo := r.AccountRepository()

	tx, err := r.Transaction()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"step":  "transaction-create",
		}).Error("Api::Accounts->Register database error.")
		render.Status(http.StatusInternalServerError)
		return
	}

	var host string

	if len(data.Account.Host) == 0 {
		host, err = generateHost(r, config)
	} else {
		var account *doorbot.Account
		account, err = repo.FindByHost(r.DB(), data.Account.Host)
		if account != nil {
			tx.Rollback()

			render.Status(http.StatusConflict)
			return
		}

		host = data.Account.Host
	}

	if err != nil {
		tx.Rollback()
		render.Status(http.StatusInternalServerError)
		return
	}

	if len(host) == 0 {
		log.WithFields(log.Fields{
			"host": host,
			"step": "host-generation",
		}).Error("Api::Accounts->Register Unable to set a hostname.")

		tx.Rollback()
		render.Status(http.StatusServiceUnavailable)
		return
	}

	// Create the account instance
	account := &doorbot.Account{
		Name:         data.Account.Name,
		ContactName:  data.Account.ContactName,
		ContactEmail: data.Account.ContactEmail,

		//TODO append the doorbot production domain
		Host: host,

		IsEnabled: true,
	}

	err = repo.Create(tx, account)

	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"error": err,
			"host":  host,
			"step":  "account-create",
		}).Error("Api::Accounts->Register database error.")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	// Update the repositories account scopes to the one we just created.
	r.SetAccountScope(account.ID)

	// We need to create a person with a password to let them log in on the dashboard.
	person := &doorbot.Person{
		Name:        data.Account.ContactName,
		Email:       data.Account.ContactEmail,
		AccountType: doorbot.AccountOwner,
	}

	ar := r.PersonRepository()

	err = ar.Create(tx, person)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"host":  host,
			"email": person.Email,
			"step":  "person-create",
		}).Error("Api::Accounts->Register database error.")

		tx.Rollback()

		render.Status(http.StatusInternalServerError)
		return
	}

	// Generate a random password
	password := security.RandomPassword(8)
	hash, err := security.PasswordCrypt([]byte(password))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"host":  host,
			"email": person.Email,
			"step":  "person-password",
		}).Error("Api::Accounts->Register password generation error.")

		tx.Rollback()
		render.Status(http.StatusInternalServerError)
		return
	}

	// Create a new authentication record for the user
	authr := r.AuthenticationRepository()
	authentication := &doorbot.Authentication{
		PersonID:   person.ID,
		ProviderID: auth.ProviderPassword,
		Token:      string(hash),
	}

	err = authr.Create(tx, authentication)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"host":  host,
			"email": person.Email,
			"step":  "person-authentication",
		}).Error("Api::Accounts->Register database error.")

		tx.Rollback()

		render.Status(http.StatusInternalServerError)
		return
	}
	err = tx.Commit()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"host":  host,
			"email": person.Email,
			"step":  "transaction-commit",
		}).Error("Api::Accounts->Register database error.")

		tx.Rollback()

		render.Status(http.StatusInternalServerError)
		return
	}

	//TODO Send an email to the user.

	log.WithFields(log.Fields{
		"account_id":   account.ID,
		"account_host": account.Host,
		"person_id":    person.ID,
	}).Info("Account created")

	n.AccountCreated(account, person, password)

	render.JSON(http.StatusCreated, AccountViewModel{Account: account})
}

// generateHost generates a host that isn't yet taken.
func generateHost(r doorbot.Repositories, config *doorbot.DoorbotConfig) (string, error) {
	repo := r.AccountRepository()

	// Generate a random valid host name
	for i := 0; i < 10; i++ {
		tmp := strings.ToLower(security.RandomPassword(10))

		exists, err := repo.FindByHost(r.DB(), tmp)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"host":  tmp,
			}).Error("Api::Accounts->Register database error.")

			return "", err
		}

		if exists == nil {
			return tmp, nil
		}
	}

	return "", nil
}

// Post create a new account ( using the admin panel )
func Post(render render.Render, r doorbot.Repositories, vm AccountViewModel, administrator *doorbot.Administrator) {
	repo := r.AccountRepository()

	exists, err := repo.FindByHost(r.DB(), vm.Account.Host)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"host":  vm.Account.Host,
		}).Error("Api::Accounts->Post database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if exists != nil {
		log.WithFields(log.Fields{
			"host": vm.Account.Host,
		}).Warn("Api::Accounts->Post Host already registered")

		render.JSON(http.StatusConflict, doorbot.NewConflictErrorResponse([]string{"The specified host is already registered."}))
		return
	}

	err = repo.Create(r.DB(), vm.Account)

	if err != nil {

		log.WithFields(log.Fields{
			"error": err,
			"host":  vm.Account.Host,
		}).Error("Api::Accounts->Post database error.")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"administrator_id": administrator.ID,
		"account_id":       vm.Account.ID,
		"host":             vm.Account.Host,
	}).Info("Account created by administrator")

	render.JSON(http.StatusCreated, vm)
}

// Put updates an account ( using the admin panel )
func Put(render render.Render, a *doorbot.Account, r doorbot.Repositories, vm AccountViewModel, session *auth.Authorization) {
	repo := r.AccountRepository()

	switch session.Type {
	case auth.AuthorizationAdministrator:
		// ok
	case auth.AuthorizationPerson:
		if !session.Person.IsAccountManager() {

			log.WithFields(log.Fields{
				"person_id":  session.Person.ID,
				"account_id": a.ID,
			}).Warn("Api::Accounts->Put unauthorized user attempted to modify account.")

			render.Status(http.StatusForbidden)
			return
		}
	}

	a.Name = vm.Account.Name

	a.BridgeHubEnabled = vm.Account.BridgeHubEnabled
	a.BridgeHubURL = vm.Account.BridgeHubURL
	a.BridgeHubToken = vm.Account.BridgeHubToken

	a.BridgeHipChatEnabled = vm.Account.BridgeHipChatEnabled
	a.BridgeHipChatToken = vm.Account.BridgeHipChatToken

	a.BridgeSlackEnabled = vm.Account.BridgeSlackEnabled
	a.BridgeSlackToken = vm.Account.BridgeSlackToken

	a.NotificationsEmailMessageTemplate = vm.Account.NotificationsEmailMessageTemplate
	a.NotificationsSMSMessageTemplate = vm.Account.NotificationsSMSMessageTemplate

	a.NotificationsEnabled = vm.Account.NotificationsEnabled

	a.NotificationsMailgunEnabled = vm.Account.NotificationsMailgunEnabled

	a.NotificationsPostmarkEnabled = vm.Account.NotificationsPostmarkEnabled


	a.NotificationsMailgunEnabled = vm.Account.NotificationsMailgunEnabled

	a.NotificationsNexmoEnabled = vm.Account.NotificationsNexmoEnabled
	a.NotificationsNexmoToken = vm.Account.NotificationsNexmoToken

	a.NotificationsSlackEnabled = vm.Account.NotificationsSlackEnabled
	a.NotificationsSlackToken = vm.Account.NotificationsSlackToken
	
	a.NotificationsTwilioEnabled = vm.Account.NotificationsTwilioEnabled


	_, err := repo.Update(r.DB(), a)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": a.ID,
		}).Error("Api::Accounts->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	vm.Account = a

	render.JSON(http.StatusOK, vm)
}

// Delete an account ( admin panel )
func Delete(render render.Render, r doorbot.Repositories, params martini.Params, administrator *doorbot.Administrator) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.AccountRepository()

	account, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":            err.Error(),
			"account_id":       account.ID,
			"administrator_id": administrator.ID,
		}).Error("Api::Accounts->Delete database find error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if account == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified account does not exists."}))
		return
	}

	_, err = repo.Delete(r.DB(), account)

	if err != nil {
		log.WithFields(log.Fields{
			"error":            err.Error(),
			"administrator_id": administrator.ID,
			"account_id":       account.ID,
		}).Error("Api::Accounts->Delete database delete error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"administrator_id": administrator.ID,
		"account_id":       account.ID,
	}).Info("Api::Accounts->Delete account deleted by administrator.")

	render.Status(http.StatusNoContent)
}
