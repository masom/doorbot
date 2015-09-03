// Package api represents the API layer
package api

import (
	"github.com/masom/doorbot/doorbot/api/accounts"
	"github.com/masom/doorbot/doorbot/api/auth"
	"github.com/masom/doorbot/doorbot/api/devices"
	"github.com/masom/doorbot/doorbot/api/doors"
	"github.com/masom/doorbot/doorbot/api/notifications"
	"github.com/masom/doorbot/doorbot/api/people"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
)

// BindRoutes binds API routes.
func BindRoutes(r martini.Router) {

	r.NotFound(func(req *http.Request, r render.Render) {
		r.JSON(http.StatusNotFound, []string{"The requested url does not exists."})
	})

	r.Get("/api/status", Status)

	// Registration
	r.Group("/api", func(r martini.Router) {
		r.Group("/accounts", func(r martini.Router) {
			r.Post("/register", binding.Bind(accounts.RegisterViewModel{}), NotificatorHandler(), accounts.Register)
		})
	}, GatekeeperRouteHandler())

	// DASHBOARD GATEKEEPER
	r.Group("/api", func(r martini.Router) {
		r.Group("/auth", func(r martini.Router) {
			r.Post("/password", binding.Bind(auth.PasswordRequest{}), auth.Password)
		})
	}, GatekeeperRouteHandler(), AccountScopeHandler(), RepositoryScopeHandler())

	//ADMIN
	r.Group("/api", func(r martini.Router) {
		r.Get("/accounts", accounts.Index)
		r.Post("/accounts", binding.Bind(accounts.AccountViewModel{}), accounts.Post)
		r.Delete("/:id", accounts.Delete)
	}, AdminRouteHandler())

	//Owner / Manager
	r.Group("/api", func(r martini.Router) {
		r.Group("/accounts", func(r martini.Router) {
			r.Get("/:id", accounts.Get)
			r.Put("/:id", binding.Bind(accounts.AccountViewModel{}), accounts.Put)
		})

		r.Get("/devices", devices.Index)
		r.Post("/devices", binding.Bind(devices.DeviceViewModel{}), devices.Post)
		r.Group("/devices", func(r martini.Router) {
			r.Get("/:id", devices.Get)
			r.Put("/:id", binding.Bind(devices.DeviceViewModel{}), devices.Put)
			r.Post("/register", binding.Bind(devices.DeviceViewModel{}), devices.Post)
			r.Delete("/:id", devices.Delete)
		})

		r.Get("/doors", doors.Index)
		r.Post("/doors", binding.Bind(doors.DoorViewModel{}), doors.Post)
		r.Group("/doors", func(r martini.Router) {
			r.Get("/:id", doors.Get)
			r.Put("/:id", binding.Bind(doors.DoorViewModel{}), doors.Put)
			r.Delete("/:id", doors.Delete)
		})

		r.Post("/people", binding.Bind(people.PersonViewModel{}), people.Post)
		r.Group("/people", func(r martini.Router) {
			r.Post("/sync", BridgeHandler(), people.Sync)
			r.Delete("/:id", people.Delete)
		})

	}, AccountScopeHandler(), RepositoryScopeHandler(), SecuredRouteHandler(), ManagerRestrictedRouteHandler())

	//Account
	r.Group("/api", func(r martini.Router) {
		r.Group("/accounts", func(r martini.Router) {
			r.Get("/:id", accounts.Get)
		})

		r.Post("/notifications", NotificatorHandler(), binding.Bind(notifications.ViewModel{}), notifications.Notify)

		r.Get("/people", people.Index)
		r.Group("/people", func(r martini.Router) {
			r.Get("/:id", people.Get)
			r.Put("/:id", binding.Bind(people.PersonViewModel{}), people.Put)
		})

	}, AccountScopeHandler(), RepositoryScopeHandler(), SecuredRouteHandler())
}
