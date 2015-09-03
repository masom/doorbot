package postmark

import(
	log "github.com/Sirupsen/logrus"
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/rendering"
	"github.com/gcmurphy/postmark"
	"fmt"
)

type Config struct {
	Token string
}

type Postmark struct {
	Account *doorbot.Account
	Token string
}

func New(a *doorbot.Account, c Config) *Postmark {
	return &Postmark{
		Account: a,
		Token: c.Token,
	}
}

// Name returns the notifier name
func (p *Postmark) Name() string {
	return "Postmark"
}

func (p *Postmark) AccountCreated(person *doorbot.Person, password string) error {
	log.WithFields(log.Fields{
		"account_id":   p.Account.ID,
		"person_id":    person.ID,
		"person_email": person.Email,
	}).Info("Notificator::Mailgun->AccountCreated")

	message := &postmark.Message{
		From:    "martin@doorbot.co", //TODO make this a configuration value
		To:      person.Email,
		Subject: "Doorbot - Account Created",
		TextBody: fmt.Sprintf(
			"Welcome %s,\n\nYou can log in on the dashboard using this temporary password: %s\n\n\nAccount: %d\nTemporary Host: %s\n Email: %s\nPassword: %s\n\n- Doorbot",
			person.Name, password, p.Account.ID, p.Account.Host, person.Email, password,
		),
	}

	pm := postmark.NewPostmark(p.Token)
	_, err := pm.Send(message)

	if err != nil {
		log.WithFields(log.Fields{
			"account_id":   p.Account.ID,
			"person_id":    person.ID,
			"person_email": person.Email,
		}).Error("Notificator::EmailNotifier->AccountCreated error")

		return err
	}

	return nil
}

// KnockKnock sends a private message to a user.
func (p *Postmark) KnockKnock(d *doorbot.Door, person *doorbot.Person) error {
	log.WithFields(log.Fields{
		"account_id": p.Account.ID,
		"person_id":  person.ID,
		"door_id":    d.ID,
	}).Info("Notificator::Postmark->KnockKnock request")

	dbb := rendering.DoorbotBar()

	renderingData := map[string]string{
		"name": person.Name,
		"door": d.Name,
	}

	//TODO HTML message template.
	message := &postmark.Message{
		From:     "martin@canvaspop.com",
		To:       person.Email,
		Subject:  dbb.Render("Doorbot - There is someone waiting at the {{door}}.", renderingData),
		TextBody: dbb.Render("Hi {{name}},\nThere is someone waiting at the {{door}}.\n\n - Doorbot", renderingData),
	}

	pm := postmark.NewPostmark(p.Token)

	_, err := pm.Send(message)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": p.Account.ID,
			"person_id":  person.ID,
			"door_id":    d.ID,
		}).Error("Notificator::Postmark->KnockKnock error")
		return err
	}

	return nil
}
