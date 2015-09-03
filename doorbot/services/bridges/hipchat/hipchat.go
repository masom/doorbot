package hipchat

import (
	"github.com/tbruyelle/hipchat-go/hipchat"
	log "github.com/Sirupsen/logrus"
	"strconv"
)


type HipChatConfig struct {
	Token string
}

type HipChat struct {
	AccountID uint
	Token string
}


type HipChatUser struct {
	ID int
	Email string
	DisplayName string
	Title string
}

func New(accountID uint, config *HipChatConfig) *HipChat {
	return &HipChat{
		AccountID: accountID,
		Token: config.Token,
	}
}


func (h *HipChat) GetUsers() ([]*HipChatUser, error) {
	log.WithFields(log.Fields{
		"account_id": h.AccountID,
		"token": h.Token,
	}).Info("Bridges::HipChat::GetUsers started")

	c := hipchat.NewClient(h.Token)

	svcUsers, _, err := c.User.List(0, 1000, false, false)
	if err != nil {
		log.WithFields(log.Fields{
			"account_id": h.AccountID,
			"error": err,
		}).Error("Bridges::HipChat::GetUsers error")
	}

	users := make([]*HipChatUser, len(svcUsers))

	for k, u := range svcUsers {
		details, _, err := c.User.View(strconv.FormatInt(int64(u.ID), 10))

		if err != nil {
			return users, err
		}

		users[k] = &HipChatUser {
			ID: u.ID,
			Email: details.Email,
			DisplayName: details.Name,
			Title: details.Title,
		}
	}

	return users, err
}
