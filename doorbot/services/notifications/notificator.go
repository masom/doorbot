package notifications

import (
	log "github.com/Sirupsen/logrus"
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/notifications/hipchat"
	"github.com/masom/doorbot/doorbot/services/notifications/mailgun"
	"github.com/masom/doorbot/doorbot/services/notifications/nexmo"
	"github.com/masom/doorbot/doorbot/services/notifications/postmark"
	"github.com/masom/doorbot/doorbot/services/notifications/slack"
	"github.com/masom/doorbot/doorbot/services/notifications/twilio"
)

type (
	// Notificator configuration
	Config struct {
		Account *doorbot.Account
		HipChat hipchat.Config
		Mailgun mailgun.Config
		Nexmo nexmo.Config
		Postmark postmark.Config
		Slack slack.Config
		Twilio twilio.Config
	}

	Notificator interface {
		AccountCreated(a *doorbot.Account, p *doorbot.Person, password string)
		KnockKnock(d *doorbot.Door, p *doorbot.Person) bool
	}

	// Notifier interface
	Notifier interface {
		KnockKnock(d *doorbot.Door, p *doorbot.Person) error
		Name() string
	}

	// Notificator decides how to notify given events
	notificator struct {
		Config Config
	}
)

// New create a new Notificator
func New(c Config) Notificator {
	return &notificator{
		Config: c,
	}
}

// AccountCreated sends an email to Doorbot
func (n *notificator) AccountCreated(a *doorbot.Account, p *doorbot.Person, password string) {
	en := mailgun.New(a, n.Config.Mailgun)
	go en.AccountCreated(p, password)
}

// KnockKnock sends a notification to a user that someone is looking for them at a certain door
func (n *notificator) KnockKnock(d *doorbot.Door, p *doorbot.Person) bool {
	channels := n.channels(p)

	if len(channels) == 0 {
		log.WithFields(log.Fields{
			"account_id": n.Config.Account.ID,
			"person_id": p.ID,
			"door_id": d.ID,
		}).Info("Notificator::KnockKnock no channels")
		return false
	}

	go func(channels []Notifier){


		for _, channel := range channels {
			err := channel.KnockKnock(d, p)

			if err == nil {
				log.WithFields(log.Fields{
					"account_id": n.Config.Account.ID,
					"person_id": p.ID,
					"door_id": d.ID,
					"channel":  channel.Name(),
					"error": err,
				}).Error("Notificator::KnockKnock error")

				continue
			}

			// TODO check if the user want to be contacted on multiple channels
			log.WithFields(log.Fields{
				"account_id": n.Config.Account.ID,
				"person_id": p.ID,
				"door_id": d.ID,
				"channel":  channel.Name(),
			}).Info("Notificator::KnockKnock delivered")

			return
		}
	}(channels)

	return true
}

// channels builds a list of channels the account + user have enabled.
func (n *notificator) channels(p *doorbot.Person) []Notifier {

	var notifiers []Notifier

	if p.NotificationsChatEnabled {
		// HipChat
		if n.Config.Account.NotificationsHipChatEnabled {
			notifiers = append(notifiers, hipchat.New(n.Config.Account, n.Config.HipChat))
		}

		// Slack
		if n.Config.Account.NotificationsSlackEnabled {
			notifiers = append(notifiers, slack.New(n.Config.Account, n.Config.Slack))
		}
	}

	if p.NotificationsSMSEnabled {
		// Nexmo
		if n.Config.Account.NotificationsNexmoEnabled {
			if  len(p.PhoneNumber) > 6 {
				notifiers = append(notifiers, nexmo.New(n.Config.Account, n.Config.Nexmo))
			}
		}

		// Twilio
		if n.Config.Account.NotificationsTwilioEnabled {
			if len(p.PhoneNumber) > 6 {
				notifiers = append(notifiers, twilio.New(n.Config.Account, n.Config.Twilio))
			}
		}
	}

	if p.NotificationsEmailEnabled {
		// Mailgun
		if n.Config.Account.NotificationsMailgunEnabled {
			notifiers = append(notifiers, mailgun.New(n.Config.Account, n.Config.Mailgun))
		}

		// Postmark
		if n.Config.Account.NotificationsPostmarkEnabled {
			notifiers = append(notifiers, postmark.New(n.Config.Account, n.Config.Postmark))
		}
	}

	return notifiers
}
