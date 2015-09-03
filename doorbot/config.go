package doorbot

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
)

// DatabaseConfig holds database configuration values
type DatabaseConfig struct {
	Trace bool
	URL   string
}

// ServerConfig holds http server configuration values
type ServerConfig struct {
	Port uint
}

// DoorbotConfig holds doorbot configuration values
type DoorbotConfig struct {
	Debug   bool
	Domains map[string]interface{}

	Server   ServerConfig
	Database DatabaseConfig

	// base domain name for user accounts ex: [name].doorbot.com
	UserAccountsDomain string
}

// EnvConfig holds environment-provided configuration values
type EnvConfig struct {
	Debug   bool
	Domains string

	NotificatorEmailEnabled  bool
	NotificatorPostmarkToken string
	NotificatorEmailFrom     string

	NotificatorSmsEnabled        bool
	NotificatorTwilioId          string
	NotificatorTwilioToken       string
	NotificatorTwilioPhoneNumber string

	ServerPort    uint
	DatabaseTrace bool
	DatabaseUrl   string

	UserAccountsDomain string
}

// HerokuConfig holds configuration data provided within the Heroku environment.
type HerokuConfig struct {
	Debug       bool
	DatabaseUrl string

	NotificatorEmailEnabled  bool
	NotificatorPostmarkToken string

	NotificatorSmsEnabled        bool
	NotificatorTwilioId          string
	NotificatorTwilioToken       string
	NotificatorTwilioPhoneNumber string

	Port    uint
	Domains string

	UserAccountsDomain string
}

// ParseDomains pasrses the provided domains and set them on the configuration instance.
func (c *DoorbotConfig) ParseDomains(d string) []string {
	domains := strings.Split(d, ",")

	log.Printf(fmt.Sprintf("DoorbotConfig::ParseDomains registering the following domains as site: %s", d))

	c.Domains = map[string]interface{}{}
	for _, domain := range domains {
		c.Domains[domain] = true
	}

	return domains
}
