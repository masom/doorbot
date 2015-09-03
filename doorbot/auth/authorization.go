// Package auth holds authentication related methods and structures
package auth

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/security"
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	// AuthorizationGatekeeper sets the authorization type to gatekeeper ( the dashboard login screen )
	AuthorizationGatekeeper = "gatekeeper"
	// AuthorizationPerson sets the authorization type to person
	AuthorizationPerson = "person"
	// AuthorizationDevice sets the authorization type to device
	AuthorizationDevice = "device"
	// AuthorizationAdministrator sets the authorization type to administrator
	AuthorizationAdministrator = "administrator"
)

// Authorization holds authorization state
type Authorization struct {
	Type          string
	Administrator *doorbot.Administrator
	Person        *doorbot.Person
	Device        *doorbot.Device
	Policy        *security.Policy
}

// AuthenticateAdministrator wraps the authentication logic for administrators.
func AuthenticateAdministrator(r doorbot.Repositories, token string) (*doorbot.Administrator, error) {

	var administrator *doorbot.Administrator

	ar := r.AdministratorRepository()
	aar := r.AdministratorAuthenticationRepository()

	authentication, err := aar.FindByProviderIDAndToken(r.DB(), ProviderAPIToken, token)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Auth::Authorization::AuthenticationAdministrator database error")

		return administrator, errors.New("Not authorized")
	}

	if authentication == nil {
		log.WithFields(log.Fields{
			"token": token,
		}).Warn("Doorbot::Authorization::AuthenticateAdministrator token not found")

		return administrator, errors.New("Not authorized")
	}

	administrator, err = ar.Find(r.DB(), authentication.AdministratorID)

	if err != nil {
		log.Println(err)
		return administrator, errors.New("Doorbot::AuthenticateAdministrator ")
	}

	return administrator, nil
}

// AuthenticateDevice wraps the authentication logic for devices
func AuthenticateDevice(r doorbot.Repositories, token string) (*doorbot.Device, error) {
	dr := r.DeviceRepository()
	device, err := dr.FindByToken(r.DB(), token)

	if err != nil {
		log.WithFields(log.Fields{
			"error:":       err,
			"device_token": token,
			"step":         "device-get-by-token",
		}).Error("Api::Handlers->SecuredRouteHandler database error")

		return nil, err
	}

	if device == nil {
		log.WithFields(log.Fields{
			"account_id": r.AccountScope(),
			"device_id": device.ID,
		}).Warn("Api::Handlers->SecuredRouteHandler device not found")
		return nil, nil
	}

	if device.IsEnabled == false {
		log.WithFields(log.Fields{
			"account_id": r.AccountScope(),
			"device_id":  device.ID,
		}).Warn("Api::Handlers->SecuredRouteHandler device disabled.")

		return nil, nil
	}

	return device, nil
}

// AuthenticatePerson wraps the authentication logic for people
func AuthenticatePerson(r doorbot.Repositories, token string) (*doorbot.Person, error) {
	ar := r.AuthenticationRepository()
	pr := r.PersonRepository()

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		log.WithFields(log.Fields{
			"account_id":  r.AccountScope(),
			"token":       token,
			"parts_count": len(parts),
		}).Warn("Auth->AuthenticatePerson Invalid person token")

		return nil, nil
	}

	id, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return nil, nil
	}

	token = parts[1]
	authentication, err := ar.FindByProviderIDAndPersonIDAndToken(r.DB(), ProviderAPIToken, uint(id), token)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"step":       "authentication-get-by-person-and-token",
			"person_id":  id,
			"account_id": r.AccountScope(),
			"token":      token,
		}).Error("Auth->AuthenticatePerson database error")
		return nil, err
	}

	if authentication == nil {
		log.WithFields(log.Fields{
			"token":      token,
			"person_id":  id,
			"account_id": r.AccountScope(),
		}).Info("Auth->AuthenticatePerson authentication not found")
		return nil, nil
	}

	person, err := pr.Find(r.DB(), authentication.PersonID)

	if err != nil {
		log.WithFields(log.Fields{
			"error":             err,
			"account_id":        r.AccountScope(),
			"person_id":         id,
			"authentication_id": authentication.ID,
			"step":              "person-find",
		}).Error("Auth->AuthenticatePerson database error")

		return nil, err
	}

	if person == nil {
		log.WithFields(log.Fields{
			"token":             token,
			"person_id":         id,
			"account_id":        r.AccountScope(),
			"authentication_id": authentication.ID,
		}).Error("Auth->AuthenticatePerson person not found")

		return nil, nil
	}

	return person, err
}
