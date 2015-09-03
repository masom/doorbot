package hipchat

import (
	log "github.com/Sirupsen/logrus"
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/rendering"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

// Config holds HipChat configuration values
type Config struct {
	Token string
}


// HipChat notifier
type HipChat struct {
	Token string
	Account *doorbot.Account
}

// New creates a HipChat instance
func New(a *doorbot.Account, c Config) *HipChat {
	return &HipChat{
		Token: c.Token,
	}
}

// Name return the notifier name
func (h *HipChat) Name() string {
	return "HipChat"
}

// KnockKnock sends a private message to a user.
func (h *HipChat) KnockKnock(d *doorbot.Door, p *doorbot.Person) error {

	log.WithFields(log.Fields{
		"account_id": h.Account.ID,
		"person_id":  p.ID,
		"door_id":    d.ID,
	}).Info("Notificator::HipChat->Notify request")


	c := hipchat.NewClient(h.Token)

	dbb := rendering.DoorbotBar()

	renderingData := map[string]string{
		"name": p.Name,
		"door": d.Name,
	}

	messageRequest := &hipchat.MessageRequest{
		Message: dbb.Render("Hi {{name}},\nThere is someone waiting at the {{door}}.\n\n - Doorbot", renderingData),
		Notify: true,
	}

	response, err := c.User.Message(p.Email, messageRequest)
	if response.StatusCode != 200 {
		return err
	}

	return nil
}
