package api

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

// NewServer returns a martini instance for the API endpoints
func NewServer(c *doorbot.DoorbotConfig) *martini.Martini {
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Map(c)

	m.Use(CORSHandler())

	doorbot.MapDatabase(m.Martini, c)
	doorbot.UseRepositories(m.Martini)

	m.Use(func(req *http.Request, render render.Render) {
		if req.Method == "OPTIONS" {
			render.Status(http.StatusOK)
			return
		}
	})

	BindRoutes(m.Router)

	return m.Martini
}
