package auth

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/auth"
	"github.com/masom/doorbot/doorbot/security"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/martini-contrib/render"
	"net/http"
)

// PasswordRequest represents a password request (lol)
type PasswordRequest struct {
	Authentication PasswordAuthentication `json:"authentication" binding:"required"`
}

// PasswordAuthentication contains password authentication data
type PasswordAuthentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// APITokenResponse is a response with an API token
type APITokenResponse struct {
	Authentication TokenAuthentication `json:"authentication"`
	Person         *doorbot.Person     `json:"person"`
	Policy         *security.Policy    `json:"policy"`
	Device         *doorbot.Device     `json:"device"`
}

// TokenAuthentication contains token authentication data.
type TokenAuthentication struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

// TokenAuthenticationRequest represents a token request
type TokenRequest struct {
	Authentication TokenAuthentication `json:"authentication" binding:"required"`
}

// Password will authenticate a person using the provided email and password.
// A token will be generated if the authentication succeeds.
func Password(render render.Render, account *doorbot.Account, r doorbot.Repositories, vm PasswordRequest) {
	personRepo := r.PersonRepository()
	authRepo := r.AuthenticationRepository()

	// Find the person by email
	person, err := personRepo.FindByEmail(r.DB(), vm.Authentication.Email)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"email":      vm.Authentication.Email,
			"account_id": account.ID,
			"step":       "find-person-by-email",
		}).Error("Api::Auth->Password database error")

		render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"}))
		return
	}

	if person == nil {
		log.WithFields(log.Fields{
			"account_id": account.ID,
			"email":      vm.Authentication.Email,
			"step":       "find-person-by-email",
		}).Info("Api::Auth->Password invalid email.")

		render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"}))
		return
	}

	//Fetch the password authentication record
	authentication, err := authRepo.FindByPersonIDAndProviderID(r.DB(), person.ID, auth.ProviderPassword)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": account.ID,
			"person_id":  person.ID,
			"step":       "find-person-authentication",
		}).Error("Api::Auth->Password database error")

		render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"}))
		return
	}

	if authentication == nil {
		log.WithFields(log.Fields{
			"account_id": account.ID,
			"person_id":  person.ID,
			"step":       "find-person-authentication",
		}).Info("Api::Auth->Password no authentication")

		render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"}))
		return
	}

	//Compare the passwords
	err = security.PasswordCompare([]byte(authentication.Token), []byte(vm.Authentication.Password))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"person_id":  person.ID,
			"account_id": account.ID,
			"step":       "compare-password",
		}).Info("Api::Auth->Password compare password error")

		render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"}))
		return
	}

	// Find the first active API token for this person.
	token, err := authRepo.FindByPersonIDAndProviderID(r.DB(), person.ID, auth.ProviderAPIToken)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": account.ID,
			"person_id":  person.ID,
			"step":       "find-authentication",
		}).Error("Api::Auth->Password database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	// No active token or the person has not yet signed in. Generate a new token.
	if token == nil {
		token = &doorbot.Authentication{
			PersonID:   person.ID,
			ProviderID: auth.ProviderAPIToken,
			Token:      security.GenerateAPIToken(),
		}

		err = authRepo.Create(r.DB(), token)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"person_id":  person.ID,
				"account_id": account.ID,
				"step":       "save-authentication",
			}).Error("Api::Auth->Password database error")

			render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
			return
		}
	}

	log.WithFields(log.Fields{
		"account_id": account.ID,
		"person_id":  person.ID,
	}).Info("Api::Auth->Password user logged in")

	resp := APITokenResponse{}
	resp.Authentication.Token = fmt.Sprintf("%d.%s", person.ID, token.Token)
	resp.Person = person
	resp.Policy = getPolicyForPerson(person)

	render.JSON(http.StatusOK, resp)
}

func Token(render render.Render, account *doorbot.Account, r doorbot.Repositories, vm TokenRequest) {

	var person *doorbot.Person
	var device *doorbot.Device
	var policy *security.Policy

	switch vm.Authentication.Type {
	case auth.AuthorizationPerson:
		authRepo := r.AuthenticationRepository()
		personRepo := r.PersonRepository()

		authentication, err := authRepo.FindByProviderIDAndToken(r.DB(), auth.ProviderAPIToken, vm.Authentication.Token)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"account_id": account.ID,
				"token":      vm.Authentication.Token,
				"step":       "find-token",
			}).Error("Api::Auth->Token database error")

			render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
			return
		}

		if authentication == nil {
			log.WithFields(log.Fields{
				"account_id": account.ID,
				"token":      vm.Authentication.Token,
				"step":       "find-token",
			}).Info("Api::Auth->Token no token found.")

			render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid token"}))
			return
		}

		person, err = personRepo.Find(r.DB(), authentication.PersonID)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"account_id": account.ID,
				"token":      vm.Authentication.Token,
				"step":       "find-person",
			}).Error("Api::Auth->Token database error")

			render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
			return
		}

		if person == nil {
			log.WithFields(log.Fields{
				"account_id": account.ID,
				"token":      vm.Authentication.Token,
				"step":       "find-device",
			}).Warning("Api::Auth->Token person not found.")

			render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid token"}))
			return
		}

		policy = getPolicyForPerson(person)

		break
	case auth.AuthorizationDevice:
		deviceRepo := r.DeviceRepository()
		device, err := deviceRepo.FindByToken(r.DB(), vm.Authentication.Token)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"account_id": account.ID,
				"token":      vm.Authentication.Token,
				"step":       "find-device",
			}).Error("Api::Auth->Token database error")

			render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
			return
		}

		if device == nil {
			log.WithFields(log.Fields{
				"account_id": account.ID,
				"token":      vm.Authentication.Token,
				"step":       "find-device",
			}).Info("Api::Auth->Token no device found.")

			render.JSON(http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid token"}))
			return
		}
	}

	resp := APITokenResponse{}
	resp.Authentication.Token = vm.Authentication.Token
	resp.Person = person
	resp.Policy = policy
	resp.Device = device

	render.JSON(http.StatusOK, resp)
}

func getPolicyForPerson(person *doorbot.Person) *security.Policy {
	var policy *security.Policy

	switch person.AccountType {
	case doorbot.AccountMember:
		policy = security.NewMemberPolicy()
		break
	case doorbot.AccountManager:
		policy = security.NewManagerPolicy()
		break
	case doorbot.AccountOwner:
		policy = security.NewOwnerPolicy()
		break
	}
	return policy
}
