// +build tests

package notifications

import (
	"github.com/masom/doorbot/doorbot"
	"bitbucket.org/msamson/doorbot-api/tests"
	"net/http"
	"testing"
)

func TestNotify(t *testing.T) {
	render := new(tests.MockRender)
	db := new(tests.MockExecutor)
	repositories := new(tests.MockRepositories)

	notificator := new(tests.MockNotificator)
	peopleRepo := new(tests.MockPersonRepository)
	doorRepo := new(tests.MockDoorRepository)

	account := &doorbot.Account{
		ID: 44,
	}

	person := &doorbot.Person{
		AccountID:   44,
		ID:          45,
		Name:        "John Rambo",
		Email:       "jrambo@example.com",
		IsVisible:   true,
		IsAvailable: true,
	}

	door := &doorbot.Door{}

	notification := Notification{
		DoorID:   33,
		PersonID: 45,
	}

	repositories.On("PersonRepository").Return(peopleRepo)
	repositories.On("DoorRepository").Return(doorRepo)
	repositories.On("DB").Return(db)

	peopleRepo.On("Find", db, uint(45)).Return(person, nil)
	doorRepo.On("Find", db, uint(33)).Return(door, nil)

	notificator.On("Notify", account, door, person).Return(nil)

	vm := ViewModel{Notification: &notification}

	render.On("JSON", http.StatusAccepted, ViewModel{Notification: &notification}).Return()

	Notify(render, account, repositories, notificator, vm)

	render.Mock.AssertExpectations(t)
	peopleRepo.Mock.AssertExpectations(t)
	doorRepo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
	notificator.Mock.AssertExpectations(t)
}
