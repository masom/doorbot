package devices

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/security"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

// DevicesViewModel represents a list of devices
type DevicesViewModel struct {
	Devices []*doorbot.Device `json:"devices"`
}

// DeviceViewModel represent a device
type DeviceViewModel struct {
	Device *doorbot.Device `json:"device"`
}

// Index returns a list of devices
func Index(render render.Render, r doorbot.Repositories) {
	repo := r.DeviceRepository()

	devices, err := repo.All(r.DB())

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::Devices->Index database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	render.JSON(http.StatusOK, DevicesViewModel{Devices: devices})
}

// Get return a specific device
func Get(render render.Render, r doorbot.Repositories, params martini.Params) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.DeviceRepository()
	device, err := repo.Find(r.DB(), uint(id))

	if device == nil || err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
		}).Error("Api::Devices->Get database error")

		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified device does not exists"}))
		return
	}

	render.JSON(http.StatusOK, DeviceViewModel{Device: device})
}

// Register creates a new device
func Register(render render.Render, r doorbot.Repositories, vm DeviceViewModel) {

	repo := r.DeviceRepository()

	device, err := repo.FindByToken(r.DB(), vm.Device.Token)

	if err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"account_id":   r.AccountScope(),
			"device_token": vm.Device.Token,
			"step":         "device-find",
		}).Error("Api::Devices->Register database error")

		render.Status(http.StatusInternalServerError)
		return
	}

	if device == nil {
		render.Status(http.StatusNotFound)
		return
	}

	device.Make = vm.Device.Make
	device.DeviceID = vm.Device.DeviceID
	device.IsEnabled = true

	_, err = repo.Update(r.DB(), device)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  device.ID,
			"step":       "device-update",
		}).Error("Api::Devices->Register database error")

		render.Status(http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"device_id":  device.ID,
	}).Info("Api::Devices->Put device registered")

	render.JSON(http.StatusOK, device)
}

// Post creates a new device
func Post(render render.Render, r doorbot.Repositories, vm DeviceViewModel) {
	repo := r.DeviceRepository()

	vm.Device.Token = security.GenerateAPIToken()

	err := repo.Create(r.DB(), vm.Device)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
		}).Error("Api::Devices->Post database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"device_id":  vm.Device.ID,
	}).Info("Api::Devices->Put device created")

	render.JSON(http.StatusCreated, vm)
}

// Put updates a device
func Put(render render.Render, r doorbot.Repositories, params martini.Params, vm DeviceViewModel) {
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.DeviceRepository()
	device, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-find",
		}).Error("Api::Devices->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if device == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified device does not exists"}))
		return
	}

	device.Name = vm.Device.Name
	device.DeviceID = vm.Device.DeviceID
	device.Make = vm.Device.Make
	device.Description = vm.Device.Description
	device.IsEnabled = vm.Device.IsEnabled

	_, err = repo.Update(r.DB(), device)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-update",
		}).Error("Api::Devices->Put database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"device_id":  id,
	}).Info("Api::Devices->Put device updated")

	vm.Device = device

	render.JSON(http.StatusOK, vm)
}

// Enable a device
func Enable(render render.Render, r doorbot.Repositories, params martini.Params) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"the id must be an unsigned integer"}))
		return
	}

	repo := r.DeviceRepository()
	device, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-find",
		}).Error("Api::Devices->Enable database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if device == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified deevice does not exists."}))
		return
	}

	_, err = repo.Enable(r.DB(), device, true)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-enable",
		}).Error("Api::Devices->Enable database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"device_id":  id,
	}).Info("Api::Devices->Enable device enabled")

	render.Status(http.StatusNoContent)
}

// Disable a device
func Disable(render render.Render, r doorbot.Repositories, params martini.Params) {
	id, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"the id must be an unsigned integer"}))
		return
	}

	repo := r.DeviceRepository()
	device, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-find",
		}).Error("Api::Devices->Disable database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if device == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified deevice does not exists."}))
		return
	}

	_, err = repo.Enable(r.DB(), device, false)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-disable",
		}).Error("Api::Devices->Disable database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"device_id":  id,
	}).Info("Api::Devices->Disabled device disabled")

	render.Status(http.StatusNoContent)
}

// Delete a device
func Delete(render render.Render, r doorbot.Repositories, params martini.Params) {
	id, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		render.JSON(http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"}))
		return
	}

	repo := r.DeviceRepository()
	device, err := repo.Find(r.DB(), uint(id))
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-find",
		}).Error("Api::Devices->Delete database error")

		render.JSON(http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{}))
		return
	}

	if device == nil {
		render.JSON(http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified device does not exists"}))
		return
	}

	_, err = repo.Delete(r.DB(), device)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"account_id": r.AccountScope(),
			"device_id":  id,
			"step":       "device-delete",
		}).Error("Api::Devices->Delete database error")

		render.Status(http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"account_id": r.AccountScope(),
		"device_id":  id,
	}).Info("Api::Devices->Disabled device deleted")

	render.Status(http.StatusNoContent)
}
