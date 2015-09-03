package slack

import(
	"github.com/masom/doorbot/doorbot"
)

type Config struct {
	Token string
}

type Slack struct {
	Account *doorbot.Account
	Token string
}

func New(a *doorbot.Account, c Config) *Slack {
	return &Slack{
		Account: a,
		Token: c.Token,
	}
}

func (s *Slack) Name() string {
	return "Slack"
}

func (s *Slack) KnockKnock(d *doorbot.Door, p *doorbot.Person) error {
	return nil
}
