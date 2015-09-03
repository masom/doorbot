// Package doors wraps door-related api logic
package doors

import (
	"github.com/masom/doorbot/doorbot"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

// DoorsViewModel represents a list of doors
type DoorsViewModel struct {
	Doors []*doorbot.Door `json:"doors"`
}

// DoorViewModel represent a door
type DoorViewModel struct {
	Door *doorbot.Door `json:"door"`
}

// Index return a list of doors
func Index(render render.Render, r doorbot.Repositories) {
	repo := r.DoorRepository()

	doors, err := repo.All(r.DB())

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::Doors->Index database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	render.JSON(http.StatusOK, DoorsViewModel{Doors: doors})
}

// Get return a specific door
func Get(render render.Render, r doorbot.Repositories, params martini.Params) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.DoorRepository()
	door, err := repo.Find(r.DB(), uint(id))

	if door == nil || err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"door_id":    id,
		}).Error("Api::Doors->Get database error")

		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists"}))
		return
	}

	render.JSON(http.StatusOK, DoorViewModel{Door: door})
}

// Post creates a door
func Post(render render.Render, r doorbot.Repositories, vm DoorViewModel) {
	repo := r.DoorRepository()

	err := repo.Create(r.DB(), vm.Door)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::Doors->Post database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"door_id":    vm.Door.ID,
	}).Error("Api::Doors->Post door created")

	render.JSON(http.StatusCreated, vm)
}

// Put updates a door
func Put(render render.Render, r doorbot.Repositories, params martini.Params, vm DoorViewModel) {
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.DoorRepository()
	door, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"door_id":    id,
			"step":       "door-find",
		}).Error("Api::Doors->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if door == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists"}))
		return
	}

	door.Name = vm.Door.Name

	_, err = repo.Update(r.DB(), door)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"door_id":    id,
			"step":       "door-update",
		}).Error("Api::Doors->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"door_id":    vm.Door.ID,
	}).Error("Api::Doors->Post door updated")

	render.JSON(http.StatusOK, DoorViewModel{Door: door})
}

// Delete a door
func Delete(render render.Render, r doorbot.Repositories, params martini.Params) {
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.DoorRepository()
	door, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"door_id":    id,
			"step":       "door-find",
		}).Error("Api::Doors->Delete database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if door == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists"}))
		return
	}

	_, err = repo.Delete(r.DB(), door)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"door_id":    id,
			"step":       "door-delete",
		}).Error("Api::Doors->Delete database error")

		render.Status(http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"door_id":    door.ID,
	}).Error("Api::Doors->Post door deleted")

	render.Status(http.StatusNoContent)
}
