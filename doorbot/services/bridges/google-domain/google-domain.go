package googledomain

import (
	log "github.com/Sirupsen/logrus"
	directory "github.com/google/google-api-go-client/admin/directory_v1"
	"net/http"
	"strconv"
)

type (
	// Represents a google domain user
	GoogleUser struct {
		ID          uint
		Email       string
		DisplayName string
	}

	// Fetches google users
	Directory interface {
		GetUsers() ([]*GoogleUser, error)
	}

	userDirectory struct {
		AccountID uint
		Config    *Config
	}

	Config struct{}
)

func New(accountID uint, c *Config) Directory {
	return &userDirectory{
		AccountID: accountID,
		Config:    c,
	}
}

// Get users from a google domain
func (g *userDirectory) GetUsers() ([]*GoogleUser, error) {
	var users []*GoogleUser

	c := &http.Client{}
	svc, err := directory.New(c)

	if err != nil {
		log.WithFields(log.Fields{
			"account_id": g.AccountID,
			"error":      err.Error(),
		}).Error("Bridges::GoogleAdminDirectory::GetUsers")

		return users, err
	}

	svcUsers, err := svc.Users.List().Do()

	if err != nil {
		return users, err
	}

	users = make([]*GoogleUser, len(svcUsers.Users))

	var id uint64

	for k, u := range svcUsers.Users {
		id, err = strconv.ParseUint(u.Id, 10, 32)

		if err != nil {
			log.WithFields(log.Fields{
				"error":          err.Error(),
				"account_id":     g.AccountID,
				"google_user_id": u.Id,
			}).Error("Bridges::GoogleAdminDirectory::GetUsers invalid user id.")

			return users, err
		}

		users[k] = &GoogleUser{
			ID:          uint(id),
			Email:       u.PrimaryEmail,
			DisplayName: u.Name.FullName,
		}
	}

	log.WithFields(log.Fields{
		"account_id": g.AccountID,
		"users":      len(users),
	}).Info("Bridges::Hub::GetUsers Users received")

	return users, nil
}
