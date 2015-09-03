package api

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/notifications"
	"github.com/masom/doorbot/doorbot/auth"
	"github.com/masom/doorbot/doorbot/security"
	"github.com/masom/doorbot/doorbot/services/bridges"
	"github.com/masom/doorbot/doorbot/services/bridges/hub"
	"github.com/masom/doorbot/doorbot/services/bridges/hipchat"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
)

// RepositoryScopeHandler set the AccountId on the various repositories.
func RepositoryScopeHandler() martini.Handler {

	return func(res http.ResponseWriter, req *http.Request, c martini.Context, account *doorbot.Account, r doorbot.Repositories) {
		if account == nil {
			return
		}

		r.SetAccountScope(account.ID)
	}

}

// NotificatorHandler returns a martini.Handler that map a Notification instance
func NotificatorHandler() martini.Handler {
	return func(c martini.Context, a *doorbot.Account, config *doorbot.DoorbotConfig) {
		notificatorConfig := notifications.Config{
			Account: a,
		}

		c.MapTo(notifications.New(notificatorConfig), (*notifications.Notificator)(nil))
	}
}

// BridgeHandler returns a martini.Handler that maps Bridge instance.
func BridgeHandler() martini.Handler {
	return func(c martini.Context, a *doorbot.Account) {
		config := bridges.Config{
			Hub: &hub.HubConfig{
				URL:   a.BridgeHubURL,
				Token: a.BridgeHubToken,
			},
			HipChat: &hipchat.HipChatConfig{
				Token: a.BridgeHipChatToken,
			},
		}

		c.MapTo(bridges.New(a.ID, config), (*bridges.Bridges)(nil))
	}
}

// CORSHandler returns a martini.Handler that sets CORS settings.
func CORSHandler() martini.Handler {
	return func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "host, content-type, authorization, accept")
	}
}

func parseAuthorization(req *http.Request, render render.Render) ([]string, error) {
	var result []string

	authorization := req.Header.Get("Authorization")

	if len(authorization) == 0 {
		//TODO add client ip / hostname being accessed.
		log.WithFields(log.Fields{}).Info("Api::Handlers->parseAuthorization missing authorization header")
		render.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Missing Authorization"})
		return result, errors.New("Missing authorization header")
	}

	parts := strings.Split(authorization, " ")

	if len(parts) != 2 {
		//TODO log the client ip / hostname being accessed
		log.WithFields(log.Fields{
			"authorization": authorization,
		}).Info("Api::Handlers->parseAuthorization Invalid authorization header")
		render.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid authorization header"})
		return result, errors.New("Invalid authorization header")
	}

	return parts, nil
}

// GatekeeperRouteHandler hanldes securing the register process
func GatekeeperRouteHandler() martini.Handler {
	return func(req *http.Request, render render.Render, c martini.Context) {
		// Map a nil account
		var account *doorbot.Account
		c.Map(account)

		store := &auth.Authorization{}

		parts, err := parseAuthorization(req, render)
		if err != nil {
			return
		}

		mode := parts[0]
		token := parts[1]

		switch mode {
		case "dashboard":
			if token == "gatekeeper" {
				store.Type = auth.AuthorizationGatekeeper
				return
			} else {
				log.Warn("Api::Handlers->GatekeeperRouteHandler invalid dashboard authentication")
				render.Status(http.StatusUnauthorized)
				return
			}
		case "android":
			if token == "gatekeeper" {
				store.Type = auth.AuthorizationGatekeeper
				return
			} else {
				log.Warn("Api::Handlers->GateKeeperRouteHandler invalid android authentication")
				render.Status(http.StatusUnauthorized)
				return
			}
		default:
			log.WithFields(log.Fields{
				"mode":  mode,
				"token": token,
			}).Info("Api::Handlers->GatekeeperRouteHandler invalid authorization type")
			render.Status(http.StatusForbidden)
			return
		}
	}
}

// AdminRouteHandler verifies the request comes in from an administrator
func AdminRouteHandler() martini.Handler {
	return func(render render.Render, c martini.Context, req *http.Request, r doorbot.Repositories) {
		session := &auth.Authorization{}

		parts, err := parseAuthorization(req, render)
		if err != nil {
			return
		}

		mode := parts[0]
		token := parts[1]

		if mode != "administrator" {
			render.Status(http.StatusForbidden)
			return
		}

		administrator, err := auth.AuthenticateAdministrator(r, token)

		if err != nil {
			render.Status(http.StatusInternalServerError)
			return
		}

		if administrator == nil {
			render.Status(http.StatusUnauthorized)
			return
		}

		session.Type = auth.AuthorizationAdministrator
		log.WithFields(log.Fields{
			"administrator_id": administrator.ID,
			"url":              req.URL,
		}).Info("Api::Handlers->AdminRouteHandler admin request")

		c.Map(session)
		c.Map(administrator)
	}
}

// ManagerRestrictedRouteHandler ensure the request is made on behalf of the owner / account manager
func ManagerRestrictedRouteHandler() martini.Handler {

	return func(render render.Render, c martini.Context, req *http.Request, a *doorbot.Account, session *auth.Authorization) {
		// Let adminstrators access the resources
		if session.Type == auth.AuthorizationAdministrator {
			return
		}

		// Let owner / manager access the resources
		if session.Type == auth.AuthorizationPerson && session.Person != nil {
			if session.Person.IsAccountManager() {
				return
			}
		}

		render.Status(http.StatusForbidden)
	}

}

// SecuredRouteHandler returns a martini.Handler responsible of handling authentication
func SecuredRouteHandler() martini.Handler {
	return func(render render.Render, c martini.Context, req *http.Request, r doorbot.Repositories, a *doorbot.Account) {

		if a == nil {
			log.Println("Doorbot::SecuredRouteHandler No account mapped.")
			render.Status(http.StatusForbidden)
			return
		}

		session := &auth.Authorization{}

		parts, err := parseAuthorization(req, render)
		if err != nil {
			return
		}

		mode := parts[0]
		token := parts[1]

		switch mode {
		case "administrator":
			administrator, err := auth.AuthenticateAdministrator(r, token)

			if err != nil {
				render.Status(http.StatusInternalServerError)
				return
			}

			if administrator == nil {
				log.WithFields(log.Fields{
					"administrator_token": token,
					"url": req.URL,
				}).Warn("Api::Handlers->SecuredRouteHandler administrator not found.")
				render.Status(http.StatusUnauthorized)
				return
			}

			log.WithFields(log.Fields{
				"url":              req.URL,
				"administrator_id": administrator.ID,
			}).Info("Api::Handlers->SecuredRouteHandler adminstrator request")

			c.Map(administrator)
			session.Administrator = administrator
			session.Type = auth.AuthorizationAdministrator
			session.Policy = security.NewAdministratorPolicy()

		case "device":

			device, err := auth.AuthenticateDevice(r, token)

			if err != nil {
				render.Status(http.StatusInternalServerError)
				return
			}

			if device == nil {
				render.Status(http.StatusUnauthorized)
				return
			}

			log.WithFields(log.Fields{
				"url":        req.URL,
				"device_id":  device.ID,
				"account_id": a.ID,
			}).Info("Api::Handlers->SecuredRouteHandler device request")

			session.Type = auth.AuthorizationDevice
			session.Device = device
			session.Policy = security.NewDevicePolicy()

		case "person":

			person, err := auth.AuthenticatePerson(r, token)

			if err != nil {
				render.Status(http.StatusInternalServerError)
				return
			}

			if person == nil {
				render.Status(http.StatusUnauthorized)
				return
			}

			log.WithFields(log.Fields{
				"url":        req.URL,
				"account_id": a.ID,
				"person_id":  person.ID,
			}).Info("Api::Handlers->SecuredRouteHandler person request")

			session.Person = person
			session.Type = auth.AuthorizationPerson

			switch person.AccountType {
			case doorbot.AccountOwner:
				session.Policy = security.NewOwnerPolicy()
				break
			case doorbot.AccountManager:
				session.Policy = security.NewManagerPolicy()
				break
			case doorbot.AccountMember:
				session.Policy = security.NewMemberPolicy()
				break
			}

		default:
			log.WithFields(log.Fields{
				"url":        req.URL,
				"account_id": a.ID,
			}).Info("Api::Handlers->SecuredRouteHandler Unauthorized access (invalid auth type)")
			render.Status(http.StatusForbidden)
			return
		}

		c.Map(session)
	}
}

// AccountScopeHandler setup the account scope using the requested host
func AccountScopeHandler() martini.Handler {

	return func(render render.Render, c martini.Context, req *http.Request, res http.ResponseWriter, config *doorbot.DoorbotConfig, r doorbot.Repositories) {
		var account *doorbot.Account

		// Get the request host OR fallback to the host header.
		host := req.Host

		if len(host) == 0 {
			host = req.Header.Get("Host")
		}

		if len(host) == 0 {
			render.Status(http.StatusBadRequest)
			return
		}

		// Strip the :port section of the host
		if i := strings.Index(host, ":"); i >= 0 {
			host = host[:i]
		}

		// Get the subdomain ( test.doorbot.com  -> test )
		host_parts := strings.Split(host, ".")

		if len(host_parts) < 3 {
			render.Status(http.StatusForbidden)
			return
		}

		host = strings.TrimSpace(host_parts[0])

		ar := r.AccountRepository()
		account, err := ar.FindByHost(r.DB(), host)
		if err != nil {
			log.WithFields(log.Fields{
				"host": host,
				"step": "account-find-by-host",
			}).Error("Api::Handlers->AccountScopeHandler database error")

			panic(err)
		}

		// Redirect to the home page if the account does not exists.
		if account == nil {
			log.WithFields(log.Fields{
				"host": host,
			}).Info("Api::Handlers->AccountScopeHandler account not found")

			http.NotFound(res, req)
			return
		}

		c.Map(account)
	}

}
