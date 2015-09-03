package notifications

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/services/notifications"
	log "github.com/Sirupsen/logrus"
	"github.com/martini-contrib/render"
	"net/http"
)

// ViewModel wraps requests and responses.
type ViewModel struct {
	Notification *Notification `json:"notification" binding:"required"`
}

// Notification represents a notification request.
type Notification struct {
	DoorID   uint `json:"door_id" binding:"required"`
	PersonID uint `json:"person_id" binding:"required"`
}

// Notify someone their presence is needed at a given door.
func Notify(render render.Render, account *doorbot.Account, r doorbot.Repositories, notificator notifications.Notificator, vm ViewModel) {

	notification := vm.Notification

	log.WithFields(log.Fields{
		"account_id": account.ID,
		"person_id":  notification.PersonID,
		"door_id":    notification.DoorID,
	}).Info("Api::Notifications->Notify request")

	peopleRepo := r.PersonRepository()
	doorRepo := r.DoorRepository()

	person, err := peopleRepo.Find(r.DB(), notification.PersonID)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": account.ID,
			"person_id":  notification.PersonID,
			"door_id":    notification.DoorID,
			"step":       "person-find",
		}).Error("Api::Notifications->Notify database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if person == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists."}))
		log.WithFields(log.Fields{
			"account_id": account.ID,
			"person_id":  notification.PersonID,
			"door_id":    notification.DoorID,
		}).Info("Api::Notifications->Notify person not found")
		return
	}

	if !person.IsVisible || !person.IsAvailable {
		log.WithFields(log.Fields{
			"account_id":          account.ID,
			"person_id":           notification.PersonID,
			"door_id":             notification.DoorID,
			"person_is_visible":   person.IsVisible,
			"person_is_available": person.IsAvailable,
		}).Info("Api::Notifications->Notify person is not available/visible")

		//TODO Would there be a better status code?
		render.JSON(http.StatusForbidden, doorbot.NewForbiddenErrorResponse([]string{"The specified user is currently not available."}))
		return
	}

	//TODO Infer from the  device token?
	door, err := doorRepo.Find(r.DB(), notification.DoorID)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": account.ID,
			"person_id":  notification.PersonID,
			"door_id":    notification.DoorID,
			"step":       "door-find",
		}).Error("Api::Notifications->Notify database error")
		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if door == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists."}))
		log.WithFields(log.Fields{
			"account_id": account.ID,
			"person_id":  person.ID,
			"door_id":    notification.DoorID,
		}).Info("Api::Notifications->Notify door not found")
		return
	}

	notificator.KnockKnock(door, person)

	render.JSON(http.StatusAccepted, ViewModel{Notification: notification})
}
