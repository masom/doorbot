package mailgun

import(
	log "github.com/Sirupsen/logrus"
	"github.com/masom/doorbot/doorbot"
)


type Config struct {
}

type Mailgun struct {
	Account *doorbot.Account
}

func New(a *doorbot.Account, c Config) *Mailgun {
	return &Mailgun {
		Account: a,
	}
}

func (m *Mailgun) Name() string {
	return "Mailgun"
}

func (m *Mailgun) AccountCreated(p *doorbot.Person, password string ) error {
	log.WithFields(log.Fields{
		"account_id":   m.Account.ID,
		"person_id":    p.ID,
		"person_email": p.Email,
	}).Info("Notificator::Mailgun->AccountCreated")

	/**
	message := &postmark.Message{
		From:    "martin@doorbot.co", //TODO make this a configuration value
		To:      p.Email,
		Subject: "Doorbot - Account Created",
		TextBody: fmt.Sprintf(
			"Welcome %s,\n\nYou can log in on the dashboard using this temporary password: %s\n\n\nAccount: %d\nTemporary Host: %s\n Email: %s\nPassword: %s\n\n- Doorbot",
			p.Name, password, a.ID, a.Host, p.Email, password,
		),
	}
	*/
	return nil
}

func (m *Mailgun) KnockKnock(d *doorbot.Door, p *doorbot.Person) error {
	return nil
}
