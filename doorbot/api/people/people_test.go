package people

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/auth"
	"github.com/masom/doorbot/doorbot/services/bridges"
	"bitbucket.org/msamson/doorbot-api/tests"
	"errors"
	"github.com/go-martini/martini"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIndex(t *testing.T) {
	people := []*doorbot.Person{
		&doorbot.Person{
			ID:        1,
			AccountID: 1,
			Name:      "A",
		},
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)

	repo.On("All", db).Return(people, nil)
	render.On("JSON", http.StatusOK, PeopleViewModel{People: people}).Return()

	Index(render, repositories, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestIndexError(t *testing.T) {
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	people := []*doorbot.Person{}
	err := errors.New("i like pasta")

	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("AccountScope").Return(uint(0))
	repositories.On("DB").Return(db)

	repo.On("All", db).Return(people, err)
	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Index(render, repositories, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	person := &doorbot.Person{
		ID:   33,
		Name: "ACME",
	}

	render.On("JSON", http.StatusOK, PersonViewModel{Person: person}).Return(nil)
	repo.On("Find", db, uint(33)).Return(person, nil)

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestGetParseIntError(t *testing.T) {
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	render := new(tests.MockRender)

	repositories := new(tests.MockRepositories)

	params := martini.Params{
		"id": "help",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestGetNotFound(t *testing.T) {
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	var person *doorbot.Person

	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	repo.On("Find", db, uint64(33)).Return(person, nil)
	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists"})).Return()

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPost(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(0))

	person := &doorbot.Person{
		Name: "ACME",
	}

	repo.On("Create", db, person).Return(nil)

	render.On("JSON", http.StatusCreated, PersonViewModel{Person: person}).Return()

	Post(render, repositories, PersonViewModel{Person: person})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPostCreateError(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	person := &doorbot.Person{
		Name: "ACME",
	}

	repo.On("Create", db, person).Return(errors.New("errooor"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Post(render, repositories, PersonViewModel{Person: person})

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPut(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "5555",
	}

	postPerson := &doorbot.Person{
		Name: "Romanian Landlords",
	}

	repoPerson := &doorbot.Person{
		ID:          5555,
		Name:        "ACME",
		AccountType: doorbot.AccountOwner,
	}

	session := &auth.Authorization{
		Type:   auth.AuthorizationPerson,
		Person: repoPerson,
	}

	repo.On("Find", db, uint(5555)).Return(repoPerson, nil)
	repo.On("Update", db, repoPerson).Return(true, nil)

	render.On("JSON", http.StatusOK, PersonViewModel{Person: repoPerson}).Return()

	Put(render, repositories, params, PersonViewModel{Person: postPerson}, session)

	assert.Equal(t, "Romanian Landlords", repoPerson.Name)

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

	postPerson := &doorbot.Person{
		Name: "Chicken Nick",
	}

	session := &auth.Authorization{}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Put(render, repositories, params, PersonViewModel{postPerson}, session)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutNotFound(t *testing.T) {
	var person *doorbot.Person

	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "44",
	}

	postPerson := &doorbot.Person{
		Name: "Chicken Nick",
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			ID:          3,
			AccountType: doorbot.AccountManager,
		},
	}

	repo.On("Find", db, uint(44)).Return(person, nil)
	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists"})).Return()

	Put(render, repositories, params, PersonViewModel{postPerson}, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPutFailed(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "5555",
	}

	postPerson := &doorbot.Person{
		Name: "Romanian Landlords",
	}

	repoPerson := &doorbot.Person{
		ID:          5555,
		Name:        "ACME",
		AccountType: doorbot.AccountOwner,
	}

	session := &auth.Authorization{
		Type:   auth.AuthorizationPerson,
		Person: repoPerson,
	}

	repo.On("Find", db, uint(5555)).Return(repoPerson, nil)
	repo.On("Update", db, repoPerson).Return(false, errors.New("failed"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Put(render, repositories, params, PersonViewModel{Person: postPerson}, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "33",
	}

	person := &doorbot.Person{
		ID:   33,
		Name: "ACME",
	}

	account := &doorbot.Account{}
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	repo.On("Find", db, uint(33)).Return(person, nil)
	repo.On("Delete", db, person).Return(true, nil)

	render.On("Status", http.StatusNoContent).Return()

	Delete(render, repositories, params, account, session)

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

	account := &doorbot.Account{}
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Delete(render, repositories, params, account, session)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDeleteNotFound(t *testing.T) {
	var person *doorbot.Person

	repo := new(tests.MockPersonRepository)
	render := new(tests.MockRender)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "44",
	}

	account := &doorbot.Account{}
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	repo.On("Find", db, uint(44)).Return(person, nil)

	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified person does not exists"})).Return()

	Delete(render, repositories, params, account, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestDeleteFailed(t *testing.T) {
	repo := new(tests.MockPersonRepository)
	render := new(tests.MockRender)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	params := martini.Params{
		"id": "55",
	}

	person := &doorbot.Person{
		ID:   55,
		Name: "ACME",
	}

	account := &doorbot.Account{}
	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	repo.On("Find", db, uint(55)).Return(person, nil)
	repo.On("Delete", db, person).Return(false, errors.New("error"))

	render.On("Status", http.StatusInternalServerError).Return()

	Delete(render, repositories, params, account, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestSync(t *testing.T) {
	personRepo := new(tests.MockPersonRepository)
	bridgeUserRepo := new(tests.MockBridgeUserRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(personRepo)
	repositories.On("BridgeUserRepository").Return(bridgeUserRepo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(uint(1))

	render := new(tests.MockRender)

	person := &doorbot.Person{
		ID:    34,
		Name:  "Joe",
		Email: "joe@example.com",
	}
	rambo := &doorbot.Person{
		Name:  "Rambo",
		Email: "rambo@example.com",
	}

	bridgeUsers := make([]*doorbot.BridgeUser, 2)
	bridgeUsers[0] = &doorbot.BridgeUser{
		UserID: "34",
		Name:   "Bob",
		Email:  "joe+bob@example.com",
	}

	bridgeUsers[1] = &doorbot.BridgeUser{
		UserID: "35",
		Name:   "Rambo",
		Email:  "rambo@example.com",
	}

	registeredUsers := make([]*doorbot.BridgeUser, 1)
	registeredUsers[0] = &doorbot.BridgeUser{
		PersonID: 34,
		UserID:   "34",
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			AccountType: doorbot.AccountOwner,
		},
	}

	transaction := new(tests.MockTransaction)

	repositories.On("Transaction").Return(transaction, nil)
	transaction.On("Commit").Return(nil)

	bs := new(tests.MockBridges)
	bs.On("GetUsers", bridges.BridgeHub).Return(bridgeUsers, nil)

	personRepo.On("Find", transaction, uint(34)).Return(person, nil)
	personRepo.On("Create", transaction, rambo).Return(nil)
	personRepo.On("Update", transaction, person).Return(true, nil)

	bridgeUserRepo.On("FindByBridgeID", db, uint(1)).Return(registeredUsers, nil)
	bridgeUserRepo.On("Create", transaction, bridgeUsers[1]).Return(nil)
	render.On("Status", http.StatusNoContent).Return()

	Sync(render, repositories, bs, session)

	render.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
	personRepo.Mock.AssertExpectations(t)
	bridgeUserRepo.Mock.AssertExpectations(t)
	bs.Mock.AssertExpectations(t)
}
