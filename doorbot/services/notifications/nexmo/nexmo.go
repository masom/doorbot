package nexmo

import(
	"github.com/masom/doorbot/doorbot"
)

type Config struct {
	Token string
}

type Nexmo struct {
	Account *doorbot.Account
	Token string
}

func New(a *doorbot.Account, c Config) *Nexmo {
	return &Nexmo{
		Account: a,
		Token: c.Token,
	}
}

func (s *Nexmo) Name() string {
	return "Nexmo"
}

func (s *Nexmo) KnockKnock(d *doorbot.Door, p *doorbot.Person) error {
	return nil
}
