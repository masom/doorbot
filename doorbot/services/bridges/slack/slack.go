package slack

import (
	"github.com/nlopes/slack"
	log "github.com/Sirupsen/logrus"
)


type SlackConfig struct {
	Token string
}

type Slack struct {
	AccountID uint
	Token string
}

type SlackUser struct {
	ID string
	Email string
	Name string
	Title string
}

func New(accountID uint, config *SlackConfig) *Slack {
	return &Slack{
		AccountID: accountID,
		Token: config.Token,
	}
}

func (s *Slack) GetUsers() ([]*SlackUser, error) {
	log.WithFields(log.Fields{
		"account_id": s.AccountID,
		"token": s.Token,
	}).Info("Bridges::Slack::GetUsers started")

	c := slack.New(s.Token)

	svcUsers, err := c.GetUsers()
	if err != nil {
		log.WithFields(log.Fields{
			"account_id": s.AccountID,
			"error": err,
		}).Error("Bridges::Slack::GetUsers error")
	}

	var users []*SlackUser

	for _, u := range svcUsers {
		if u.Deleted || u.IsBot || u.IsRestricted || u.IsUltraRestricted {
			continue
		}

		users = append(users, &SlackUser {
			ID: u.Id,
			Email: u.Profile.Email,
			Name: u.Profile.RealName,
			Title: u.Profile.Title,
		})
	}

	return users, err
}
