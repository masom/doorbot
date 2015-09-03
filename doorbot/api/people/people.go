// Package people wraps people-related api logic
package people

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/auth"
	"github.com/masom/doorbot/doorbot/services/bridges"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

// PublicPerson represents publicly visible person data
type PublicPerson struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsAvailable bool   `json:"is_available"`
}

// PeopleViewModel represents a list of people
type PeopleViewModel struct {
	People []*doorbot.Person `json:"people"`
}

// PersonViewModel represents a person
type PersonViewModel struct {
	Person *doorbot.Person `json:"person"`
}

// PublicPersonViewModel represents a publicly visibile person data
type PublicPersonViewModel struct {
	Person *PublicPerson `json:"person"`
}

// PublicPeopleViewModel represents a list of publicly visible people
type PublicPeopleViewModel struct {
	People []*PublicPerson `json:"people"`
}

// Transform a list of poeople into a public version of the data
func newPublicPeople(people []*doorbot.Person) []*PublicPerson {
	publicPeople := make([]*PublicPerson, len(people))

	for i, p := range people {
		publicPeople[i] = newPublicPerson(p)
	}

	return publicPeople
}

// Transform a person into a public version of the data
func newPublicPerson(p *doorbot.Person) *PublicPerson {
	return &PublicPerson{
		ID:          p.ID,
		Name:        p.Name,
		Email:       p.Email,
		IsAvailable: p.IsAvailable,
	}
}

// Index returns people
func Index(render render.Render, r doorbot.Repositories, session *auth.Authorization) {
	repo := r.PersonRepository()

	people, err := repo.All(r.DB())

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::People->Index database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	switch session.Type {
	case auth.AuthorizationAdministrator:
		render.JSON(http.StatusOK, PeopleViewModel{People: people})
	case auth.AuthorizationDevice:
		render.JSON(http.StatusOK, PublicPeopleViewModel{People: newPublicPeople(people)})
	case auth.AuthorizationPerson:
		if session.Person.IsAccountManager() {
			render.JSON(http.StatusOK, PeopleViewModel{People: people})
			return
		}

		render.JSON(http.StatusOK, PublicPeopleViewModel{People: newPublicPeople(people)})
	}
}

// Get a specific person
func Get(render render.Render, r doorbot.Repositories, params martini.Params, session *auth.Authorization) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.PersonRepository()
	person, err := repo.Find(r.DB(), uint(id))

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"person_id":  id,
		}).Error("Api::People->Get database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if person == nil {
		err := doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists"})
		render.JSON(http.StatusNotFound, err)
		return
	}

	switch session.Type {
	case auth.AuthorizationAdministrator:
		render.JSON(http.StatusOK, PersonViewModel{Person: person})

	case auth.AuthorizationDevice:
		render.JSON(http.StatusOK, PublicPersonViewModel{Person: newPublicPerson(person)})

	case auth.AuthorizationPerson:
		// Display detailed info if the requesting user is an account manager or it is the same person
		if session.Person.IsAccountManager() || session.Person.ID == person.ID {
			render.JSON(http.StatusOK, PersonViewModel{Person: person})
			return
		}

		render.JSON(http.StatusOK, PublicPersonViewModel{Person: newPublicPerson(person)})
	default:
		render.Status(http.StatusForbidden)
	}
}

// Post creates a new person
func Post(render render.Render, r doorbot.Repositories, vm PersonViewModel) {

	repo := r.PersonRepository()

	err := repo.Create(r.DB(), vm.Person)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::People->Post database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"person_id":  vm.Person.ID,
	}).Info("Api::People->Post person added.")

	render.JSON(http.StatusCreated, vm)
}

// Put updates a person
func Put(render render.Render, r doorbot.Repositories, params martini.Params, vm PersonViewModel, session *auth.Authorization) {
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	var logFields log.Fields
	var logMessage string
	canUpdateAccountType := false

	switch session.Type {
	case auth.AuthorizationAdministrator:
		logFields = log.Fields{
			"account_id":      r.AccountScope(),
			"person_id":       id,
			"admnistrator_id": session.Administrator.ID,
		}
		logMessage = "Api::People->Put user updated by administrator"

		canUpdateAccountType = true

	case auth.AuthorizationPerson:
		if uint(id) != session.Person.ID {
			if session.Person.IsAccountManager() {
				canUpdateAccountType = true
			} else {
				log.WithFields(log.Fields{
					"account_id":        r.AccountScope(),
					"person_id":         id,
					"request_person_id": session.Person.ID,
				}).Warn("Api::People->Delete forbidden")

				render.Status(http.StatusForbidden)
				return
			}
		}

		logFields = log.Fields{
			"account_id":        r.AccountScope(),
			"person_id":         id,
			"request_person_id": session.Person.ID,
		}

		logMessage = "Api::People->Put user updated by user"

	default:
		log.WithFields(log.Fields{
			"account_id":        r.AccountScope(),
			"person_id":         id,
			"request_person_id": session.Person.ID,
		}).Warn("Api::People->Put forbidden")

		render.Status(http.StatusForbidden)
		return
	}

	repo := r.PersonRepository()
	person, err := repo.Find(r.DB(), uint(id))

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"person_id":  id,
			"step":       "person-find",
		}).Error("Api::People->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if person == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists"}))
		return
	}

	person.Name = vm.Person.Name
	person.Email = vm.Person.Email
	person.PhoneNumber = vm.Person.PhoneNumber
	person.Title = vm.Person.Title
	person.IsVisible = vm.Person.IsVisible
	person.IsAvailable = vm.Person.IsAvailable
	person.NotificationsEnabled = vm.Person.NotificationsEnabled
	person.NotificationsAppEnabled = vm.Person.NotificationsAppEnabled
	person.NotificationsChatEnabled = vm.Person.NotificationsChatEnabled
	person.NotificationsEmailEnabled = vm.Person.NotificationsEmailEnabled
	person.NotificationsSMSEnabled = vm.Person.NotificationsSMSEnabled

	if canUpdateAccountType {
		person.AccountType = vm.Person.AccountType
	}

	_, err = repo.Update(r.DB(), person)

	if err != nil {
		log.WithFields(log.Fields{
			"error":             err,
			"account_id":        r.AccountScope(),
			"person_id":         person.ID,
			"request_person_id": session.Person.ID,
			"step":              "person-update",
		}).Error("Api::People->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	vm.Person = person

	log.WithFields(logFields).Info(logMessage)

	render.JSON(http.StatusOK, vm)
}

// Delete a person
func Delete(render render.Render, r doorbot.Repositories, params martini.Params, a *doorbot.Account, session *auth.Authorization) {
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	var logFields log.Fields
	var logMessage string

	switch session.Type {
	case auth.AuthorizationAdministrator:
		logFields = log.Fields{
			"account_id":      r.AccountScope(),
			"person_id":       id,
			"admnistrator_id": session.Administrator.ID,
		}
		logMessage = "Api::People->Delete user deleted by administrator"

	case auth.AuthorizationPerson:
		if !session.Person.IsAccountManager() {
			log.WithFields(log.Fields{
				"account_id":        r.AccountScope(),
				"person_id":         id,
				"request_person_id": session.Person.ID,
			}).Warn("Api::People->Delete forbidden")

			render.Status(http.StatusForbidden)
			return
		}

		logFields = log.Fields{
			"account_id":     r.AccountScope(),
			"person_id":      id,
			"modified_by_id": session.Person.ID,
		}

		logMessage = "Api::People->Put user deleted by user"

	default:
		log.WithFields(log.Fields{
			"account_id":        r.AccountScope(),
			"person_id":         id,
			"request_person_id": session.Person.ID,
		}).Warn("Api::People->Delete forbidden")

		render.Status(http.StatusForbidden)
		return
	}

	repo := r.PersonRepository()
	person, err := repo.Find(r.DB(), uint(id))

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"person_id":  person.ID,
			"step":       "person-find",
		}).Error("Api::People->Delete database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if person == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists"}))
		return
	}

	_, err = repo.Delete(r.DB(), person)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"person_id":  person.ID,
			"step":       "person-delete",
		}).Error("Api::People->Delete database error")

		render.Status(http.StatusInternalServerError)
		return
	}

	log.WithFields(logFields).Info(logMessage)

	render.Status(http.StatusNoContent)
}

// Sync doorbot users with a foreign data source using a bridge.
func Sync(render render.Render, r doorbot.Repositories, b bridges.Bridges, a *doorbot.Account, session *auth.Authorization) {

	var bUsers []*doorbot.BridgeUser
	var registered []*doorbot.BridgeUser
	var err error

	personRepo := r.PersonRepository()
	bUserRepo := r.BridgeUserRepository()

	var bridgesToSync = []uint{bridges.BridgeHub, bridges.BridgeHipChat}

	for _, bridgeId := range(bridgesToSync) {
		f := func() bool {
			users, err := b.GetUsers(bridgeId)

			for _, u := range(users) {
				log.WithFields(log.Fields{"user": *u}).Info("User")
			}
			if err != nil {
				log.WithFields(log.Fields{
					"error":      err,
					"account_id": a.ID,
					"bridge_id": bridgeId,
				}).Error("Api::People->Sync bridge error")

				return false
			}

			existing, err := bUserRepo.FindByBridgeID(r.DB(), bridgeId)

			if err != nil {
				log.WithFields(log.Fields{
					"error":      err,
					"account_id": r.AccountScope(),
					"step":       "bridge-user-find-by-bridge-id",
					"bridge_id": bridgeId,
				}).Error("Api::People->Sync database error")

				return false
			}

			registered = append(registered, existing...)
			bUsers = append(bUsers, users...)

			return true;
		}

		f()
	}

	tx, err := r.Transaction()

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"step":       "transaction-create",
		}).Error("Api::People->Sync database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	var buser *doorbot.BridgeUser

	for _, u := range bUsers {
		log.WithFields(log.Fields{
			"account_id":        r.AccountScope(),
			"bridge_user_id":    u.UserID,
			"bridge_user_email": u.Email,
			"bridge_user_name":  u.Name,
		}).Debug("Api::People->Sync bridge user")

		buser = findRegistered(registered, u.UserID)

		if buser != nil {

			log.WithFields(log.Fields{
				"account_id":     r.AccountScope(),
				"bridge_user_id": buser.UserID,
				"person_id":      buser.PersonID,
			}).Debug("Api::People->Sync registered user found")

			person, err := personRepo.Find(tx, buser.PersonID)
			if err != nil {
				log.WithFields(log.Fields{
					"error":          err,
					"account_id":     r.AccountScope(),
					"bridge_user_id": buser.UserID,
					"person_id":      buser.PersonID,
					"step":           "person-find-from-bridge-id",
				}).Error("Api::People->Sync database error")

				break
			}

			person.Name = u.Name
			person.Email = u.Email
			person.PhoneNumber = u.PhoneNumber

			_, err = personRepo.Update(tx, person)

			if err != nil {
				log.WithFields(log.Fields{
					"error":          err,
					"account_id":     r.AccountScope(),
					"bridge_user_id": buser.UserID,
					"person_id":      buser.PersonID,
					"step":           "person-update-from-bridge-data",
				}).Error("Api::People->Sync database error")

				break
			}

			log.WithFields(log.Fields{
				"account_id":     r.AccountScope(),
				"bridge_user_id": buser.UserID,
				"person_id":      buser.PersonID,
			}).Info("Api::People->Sync person updated from bridge data")

			continue
		}

		log.WithFields(log.Fields{
			"account_id":        r.AccountScope(),
			"bridge_user_id":    u.UserID,
			"bridge_user_email": u.Email,
			"bridge_user_name":  u.Name,
		}).Info("Api::People->Sync new bridge user")

		// User does not exists. Create them
		args := doorbot.PersonArguments{
			Name:  u.Name,
			Email: u.Email,
		}

		person := doorbot.NewPerson(args)

		err = personRepo.Create(tx, person)
		if err != nil {
			log.WithFields(log.Fields{
				"error":          err,
				"account_id":     r.AccountScope(),
				"bridge_user_id": buser.UserID,
				"step":           "person-create-from-bridge-data",
			}).Error("Api::People->Sync database error")

			break
		}

		u.PersonID = person.ID
		err = bUserRepo.Create(tx, u)
		if err != nil {
			log.WithFields(log.Fields{
				"error":          err,
				"account_id":     r.AccountScope(),
				"bridge_user_id": buser.UserID,
				"step":           "bridge-user-create",
			}).Error("Api::People->Sync database error")

			break
		}
		continue
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::People->Sync error")

		tx.Rollback()

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"step":       "transaction-commit",
		}).Error("Api::People->Sync database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
	}).Info("Api::People->Sync bridge sync completed.")

	render.Status(http.StatusNoContent)
}

// Utility function to find a registered user.
func findRegistered(users []*doorbot.BridgeUser, id string) *doorbot.BridgeUser {
	for _, u := range users {
		if u.UserID == id {
			return u
		}
	}

	return nil
}
