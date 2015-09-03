package doors

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
	doors := []*doorbot.Door{
		&doorbot.Door{
			ID:        1,
			AccountID: 1,
			Name:      "A",
		},
	}

	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)

	repo.On("All", db).Return(doors, nil)
	render.On("JSON", http.StatusOK, DoorsViewModel{Doors: doors}).Return()

	Index(render, repositories)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestIndexError(t *testing.T) {

	doors := []*doorbot.Door{}
	err := errors.New("i like pasta")

	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
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
	repo := new(tests.MockDoorRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	door := &doorbot.Door{
		ID:   33,
		Name: "ACME",
	}

	render.On("JSON", http.StatusOK, DoorViewModel{Door: door}).Return(nil)
	repo.On("Find", db, uint(33)).Return(door, nil)

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

	var door *doorbot.Door

	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "33",
	}

	repo.On("Find", db, uint(33)).Return(door, nil)
	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists"})).Return()

	Get(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPost(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("AccountScope").Return(uint(1))

	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	door := &doorbot.Door{
		Name: "ACME",
	}

	repo.On("Create", db, door).Return(nil)

	render.On("JSON", http.StatusCreated, DoorViewModel{Door: door}).Return()

	Post(render, repositories, DoorViewModel{Door: door})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPostCreateError(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	door := &doorbot.Door{
		Name: "ACME",
	}

	repo.On("Create", db, door).Return(errors.New("errooor"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Post(render, repositories, DoorViewModel{Door: door})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPut(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "5555",
	}

	postDoor := &doorbot.Door{
		Name: "Romanian Landlords",
	}

	repoDoor := &doorbot.Door{
		ID:   5555,
		Name: "ACME",
	}

	repo.On("Find", db, uint(5555)).Return(repoDoor, nil)
	repo.On("Update", db, repoDoor).Return(true, nil)

	render.On("JSON", http.StatusOK, DoorViewModel{Door: repoDoor}).Return()

	Put(render, repositories, params, DoorViewModel{Door: postDoor})

	assert.Equal(t, "Romanian Landlords", repoDoor.Name)

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

	postDoor := &doorbot.Door{
		Name: "Chicken Nick",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Put(render, repositories, params, DoorViewModel{Door: postDoor})

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutNotFound(t *testing.T) {
	var door *doorbot.Door

	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "44",
	}

	postDoor := &doorbot.Door{
		Name: "Chicken Nick",
	}

	repo.On("Find", db, uint(44)).Return(door, nil)
	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists"})).Return()

	Put(render, repositories, params, DoorViewModel{Door: postDoor})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutFailed(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "5555",
	}

	postDoor := &doorbot.Door{
		Name: "Romanian Landlords",
	}

	repoDoor := &doorbot.Door{
		ID:   5555,
		Name: "ACME",
	}

	repo.On("Find", db, uint(5555)).Return(repoDoor, nil)
	repo.On("Update", db, repoDoor).Return(false, errors.New("failed"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Put(render, repositories, params, DoorViewModel{Door: postDoor})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockDoorRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "33",
	}

	door := &doorbot.Door{
		ID:   33,
		Name: "ACME",
	}

	repo.On("Find", db, uint(33)).Return(door, nil)
	repo.On("Delete", db, door).Return(true, nil)

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
	var door *doorbot.Door

	repo := new(tests.MockDoorRepository)
	render := new(tests.MockRender)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "44",
	}

	repo.On("Find", db, uint(44)).Return(door, nil)

	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified door does not exists"})).Return()

	Delete(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDeleteFailed(t *testing.T) {
	repo := new(tests.MockDoorRepository)
	render := new(tests.MockRender)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DoorRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "55",
	}

	door := &doorbot.Door{
		ID:   55,
		Name: "ACME",
	}

	repo.On("Find", db, uint(55)).Return(door, nil)
	repo.On("Delete", db, door).Return(false, errors.New("error"))

	render.On("Status", http.StatusInternalServerError).Return()

	Delete(render, repositories, params)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}
