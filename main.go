package main

import (
	"fmt"
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/api"
	"github.com/masom/envconfig"
	"github.com/go-martini/martini"
	log "github.com/Sirupsen/logrus"
	_ "github.com/stretchr/testify/mock"
	"math/rand"
	"os"
	"time"
)

func Configure() *doorbot.DoorbotConfig {
	c := &doorbot.DoorbotConfig{}
	c.Database = doorbot.DatabaseConfig{}
	c.Server = doorbot.ServerConfig{}

	if os.Getenv("HEROKU") != "" {
		log.Info("Using HEROKU config")

		hc := &doorbot.HerokuConfig{}

		err := envconfig.Process("", hc)
		if err != nil {
			log.Panic(err)
		}

		log.Println(hc)

		c.Database.URL = hc.DatabaseUrl
		c.Server.Port = hc.Port

		c.UserAccountsDomain = hc.UserAccountsDomain

		c.ParseDomains(hc.Domains)

	} else {
		log.Info("Using ENV config.")

		ec := &doorbot.EnvConfig{}

		err := envconfig.Process("doorbot", ec)
		if err != nil {
			log.Panic(err)
		}

		log.Println(ec)
		c.ParseDomains(ec.Domains)

		c.UserAccountsDomain = ec.UserAccountsDomain

		c.Database.URL = ec.DatabaseUrl
		c.Database.Trace = ec.DatabaseTrace
		c.Server.Port = ec.ServerPort
	}

	log.Info(fmt.Sprintf("Database URL: `%s`", c.Database.URL))
	return c
}

func setupLogging() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	config := Configure()

	m := api.NewServer(config)
	m.Use(martini.Static("public", martini.StaticOptions{IndexFile: "public/index.html"}))
	m.Run()
}
