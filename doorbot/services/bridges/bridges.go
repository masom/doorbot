package bridges

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/bridges/google-domain"
	"github.com/masom/doorbot/doorbot/services/bridges/hub"
	"github.com/masom/doorbot/doorbot/services/bridges/hipchat"
	"github.com/masom/doorbot/doorbot/services/bridges/slack"
	"strconv"
)

const (
	// BridgeHub id for the hub
	BridgeHub                  uint = 1
	BridgeGoogleAdminDirectory uint = 2
	BridgeHipChat uint = 3
	BridgeSlack uint = 4
)

type (
	bridges struct {
		AccountID uint
		Config    Config
	}

	// Config holds bridge configurations
	Config struct {
		Hub                  *hub.HubConfig
		GoogleAdminDirectory *googledomain.Config
		HipChat              *hipchat.HipChatConfig
		Slack                *slack.SlackConfig
	}

	// Bridges interface
	Bridges interface {
		GetUsers(bridgeid uint) ([]*doorbot.BridgeUser, error)
	}
)

// NewBridges creates a new bridge manager
func New(accountId uint, c Config) Bridges {
	return &bridges{
		AccountID: accountId,
		Config:    c,
	}
}

// GetUsers gets users from a bridge
func (b *bridges) GetUsers(bridgeID uint) ([]*doorbot.BridgeUser, error) {
	var users []*doorbot.BridgeUser

	switch bridgeID {
	case BridgeHub:
		bridge := hub.New(b.Config.Hub)
		hubUsers, err := bridge.GetUsers()

		if err != nil {
			return users, err
		}

		users := make([]*doorbot.BridgeUser, len(hubUsers))

		var name string

		for i, hu := range hubUsers {
			if hu.DisplayName == nil {
				name = ""
			} else {
				name = *hu.DisplayName
			}

			users[i] = &doorbot.BridgeUser{
				BridgeID: BridgeHub,
				UserID:   strconv.FormatUint(uint64(hu.ID), 10),
				Name:     name,
				Email:    hu.Email,
			}
		}

		return users, err

	case BridgeHipChat:
		bridge := hipchat.New(b.AccountID, b.Config.HipChat)

		hipChatUsers, err := bridge.GetUsers()

		if err != nil {
			return users, err
		}

		users := make([]*doorbot.BridgeUser, len(hipChatUsers))

		for i, hu := range hipChatUsers {
			users[i] = &doorbot.BridgeUser{
				BridgeID: BridgeHipChat,
				UserID: strconv.FormatUint(uint64(hu.ID), 10),
				Name: hu.DisplayName,
				Email: hu.Email,
				Title: hu.Title,
			}
		}

		return users, err

	case BridgeSlack:
		bridge := slack.New(b.AccountID, b.Config.Slack)

		slackUsers, err := bridge.GetUsers()

		if err != nil {
			return users, err
		}
		
		users := make([]*doorbot.BridgeUser, len(slackUsers))

		for i, su := range slackUsers {
			users[i] = &doorbot.BridgeUser {
				BridgeID: BridgeSlack,
				UserID: su.ID,
				Name: su.Name,
				Email: su.Email,
				Title: su.Title,
			}
		}

		return users, nil

	case BridgeGoogleAdminDirectory:
		bridge := googledomain.New(b.AccountID, b.Config.GoogleAdminDirectory)
		googleUsers, err := bridge.GetUsers()
		if err != nil {
			return users, err
		}

		users := make([]*doorbot.BridgeUser, len(googleUsers))

		for i, gu := range googleUsers {
			users[i] = &doorbot.BridgeUser{
				BridgeID: BridgeGoogleAdminDirectory,
				UserID:   strconv.FormatUint(uint64(gu.ID), 10),
				Name:     gu.DisplayName,
				Email:    gu.Email,
			}
		}

		return users, nil
	}

	return users, nil
}
