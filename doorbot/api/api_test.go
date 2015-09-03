package api

import (
	"github.com/masom/doorbot/doorbot"
	"bitbucket.org/msamson/doorbot-api/tests"
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func newServer() *martini.Martini {
	m := martini.Classic()

	c := &doorbot.DoorbotConfig{}

	m.Map(c)

	m.Use(render.Renderer())

	m.MapTo(new(tests.MockRepositories), (*doorbot.Repositories)(nil))

	BindRoutes(m.Router)

	return m.Martini
}

func getDependency(m *martini.Martini, i interface{}) interface{} {
	return m.Get(inject.InterfaceOf(i)).Interface()
}
