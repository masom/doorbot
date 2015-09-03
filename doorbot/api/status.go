package api

import (
	"github.com/martini-contrib/render"
	"net/http"
)

func Status(render render.Render) {

	//TODO setup a way to determine if the status should be different than 200
	render.Status(http.StatusOK)
}
