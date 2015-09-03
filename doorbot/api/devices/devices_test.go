package devices

import (
	"github.com/masom/doorbot/doorbot"
	"bitbucket.org/msamson/doorbot-api/tests"
	"errors"
	"github.com/go-martini/martini"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIndex(t *testing.T) {
	devices := []*doorbot.Device{
		&doorbot.Device{
			ID:        1,
			AccountID: 1,
			Name:      "A",
		},
	}

	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	repo.On("All", db).Return(devices, nil)
	render.On("JSON", http.StatusOK, DevicesViewModel{Devices: devices}).Return()

	Index(render, repositories)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestIndexError(t *testing.T) {

	doors := []*doorbot.Device{}
	err := errors.New("i like pasta")

	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	repo.On("All", db).Return(doors, err)
	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Index(render, repositories)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestGet(t *testing.T) {

	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)
	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	device := &doorbot.Device{
		ID:   33,
		Name: "ACME",
	}

	render.On("JSON", http.StatusOK, DeviceViewModel{Device: device}).Return(nil)
	repo.On("Find", db, uint(33)).Return(device, nil)

	Get(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestGetParseIntError(t *testing.T) {

	render := new(tests.MockRender)

	repositories := new(tests.MockRepositories)

	params := martini.Params{
		"id": "help",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Get(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestGetNotFound(t *testing.T) {

	var device *doorbot.Device

	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	params := martini.Params{
		"id": "33",
	}

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	repo.On("Find", db, uint64(33)).Return(device, nil)
	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified device does not exists"})).Return()

	Get(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPost(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	device := &doorbot.Device{
		Name: "ACME",
	}

	repo.On("Create", db, device).Return(nil)

	render.On("JSON", http.StatusCreated, DeviceViewModel{Device: device}).Return()

	Post(render, repositories, DeviceViewModel{Device: device})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPostCreateError(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	device := &doorbot.Device{
		Name: "ACME",
	}

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	repo.On("Create", db, device).Return(errors.New("errooor"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Post(render, repositories, DeviceViewModel{Device: device})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPut(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "5555",
	}

	postDevice := &doorbot.Device{
		Name: "Romanian Landlords",
	}

	repoDevice := &doorbot.Device{
		ID:   5555,
		Name: "ACME",
	}

	repo.On("Find", db, uint64(5555)).Return(repoDevice, nil)
	repo.On("Update", db, repoDevice).Return(true, nil)

	render.On("JSON", http.StatusOK, DeviceViewModel{Device: repoDevice}).Return()

	Put(render, repositories, params, DeviceViewModel{Device: postDevice})

	assert.Equal(t, "Romanian Landlords", repoDevice.Name)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutInvalidID(t *testing.T) {
	render := new(tests.MockRender)

	repositories := new(tests.MockRepositories)

	params := martini.Params{
		"id": "help",
	}

	postDevice := &doorbot.Device{
		Name: "Chicken Nick",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Put(render, repositories, params, DeviceViewModel{Device: postDevice})

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutNotFound(t *testing.T) {
	var device *doorbot.Device

	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "44",
	}

	postDevice := &doorbot.Device{
		Name: "Chicken Nick",
	}

	repo.On("Find", db, uint64(44)).Return(device, nil)
	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified device does not exists"})).Return()

	Put(render, repositories, params, DeviceViewModel{Device: postDevice})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutFailed(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "5555",
	}

	postDevice := &doorbot.Device{
		Name: "Romanian Landlords",
	}

	repoDevice := &doorbot.Device{
		ID:   5555,
		Name: "ACME",
	}

	repo.On("Find", db, uint64(5555)).Return(repoDevice, nil)
	repo.On("Update", db, repoDevice).Return(false, errors.New("failed"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Put(render, repositories, params, DeviceViewModel{Device: postDevice})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestEnable(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "4443",
	}

	device := &doorbot.Device{
		ID:        4443,
		IsEnabled: false,
	}

	repo.On("Find", db, uint64(4443)).Return(device, nil)
	repo.On("Enable", db, device, true).Return(true, nil)

	render.On("Status", http.StatusNoContent).Return()

	Enable(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestDisable(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "4443",
	}

	device := &doorbot.Device{
		ID:        4443,
		IsEnabled: false,
	}

	repo.On("Find", db, uint64(4443)).Return(device, nil)
	repo.On("Enable", db, device, false).Return(true, nil)

	render.On("Status", http.StatusNoContent).Return()

	Disable(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDeviceRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "33",
	}

	device := &doorbot.Device{
		ID:   33,
		Name: "ACME",
	}

	repo.On("Find", db, uint64(33)).Return(device, nil)
	repo.On("Delete", db, device).Return(true, nil)

	render.On("Status", http.StatusNoContent).Return()

	Delete(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDeleteInvalidID(t *testing.T) {
	render := new(tests.MockRender)

	repositories := new(tests.MockRepositories)

	params := martini.Params{
		"id": "help",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Delete(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDeleteNotFound(t *testing.T) {
	var device *doorbot.Device

	repo := new(tests.MockDeviceRepository)
	render := new(tests.MockRender)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "44",
	}

	repo.On("Find", db, uint64(44)).Return(device, nil)

	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified device does not exists"})).Return()

	Delete(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDeleteFailed(t *testing.T) {
	repo := new(tests.MockDeviceRepository)
	render := new(tests.MockRender)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DeviceRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "55",
	}

	device := &doorbot.Device{
		ID:   55,
		Name: "ACME",
	}

	repo.On("Find", db, uint64(55)).Return(device, nil)
	repo.On("Delete", db, device).Return(false, errors.New("error"))

	render.On("Status", http.StatusInternalServerError).Return()

	Delete(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}
