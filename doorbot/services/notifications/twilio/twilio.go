package twilio

import (
	log "github.com/Sirupsen/logrus"
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/rendering"
	"github.com/sfreiberg/gotwilio"
)

type Config struct {
	Token string
	PhoneNumber string
}

type Twilio struct {
	Account *doorbot.Account
	AccountSID string
	Token string
	PhoneNumber string
}


func New(a *doorbot.Account, c Config) *Twilio {
	return &Twilio{
		Account: a,
		Token: c.Token,
		PhoneNumber: c.PhoneNumber,
	}
}

func (t *Twilio) Name() string {
	return "Twilio"
}

func (t *Twilio) KnockKnock(d *doorbot.Door, p *doorbot.Person) error {
	log.WithFields(log.Fields{
		"account_id": t.Account.ID,
		"person_id":  p.ID,
		"person_phone_number": p.PhoneNumber,
		"door_id":    d.ID,
	}).Info("Notificator::Twilio->KnockKnock request")

	from := t.PhoneNumber
	to := p.PhoneNumber

	renderingData := map[string]string{
		"name": p.Name,
		"door": d.Name,
	}

	dbb := rendering.DoorbotBar()

	template := *t.Account.NotificationsSMSMessageTemplate
	if len(template) == 0 {
		template = "Hi {{name}}, there is someone waiting for your at the {{door}}."
	}

	message := dbb.Render(template, renderingData)

	twilio := gotwilio.NewTwilioClient(t.AccountSID, t.Token)
	// was response
	_, exception, err := twilio.SendSMS(from, to, message, "", "")
	if exception != nil {
		log.WithFields(log.Fields{
			"twilio_message": exception.Message,
			"twilio_code":    exception.Code,
			"account_id":     t.Account.ID,
			"person_id":      p.ID,
			"door_id":        d.ID,
			"sms_to":         to,
			"sms_from":       from,
		}).Error("Notificator::Twilio->KnockKnock twilio exception")

		return err
	}

	return nil
}
