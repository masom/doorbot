package accounts

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/auth"
	"bitbucket.org/msamson/doorbot-api/tests"
	"errors"
	"github.com/go-martini/martini"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestIndex(t *testing.T) {

	accounts := []*doorbot.Account{
		&doorbot.Account{
			ID:        1,
			Name:      "ACME",
			IsEnabled: true,
		},
	}

	admin := &doorbot.Administrator{}

	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)

	repo.On("All", db).Return(accounts, nil)

	render.On("JSON", http.StatusOK, AccountsViewModel{Accounts: accounts}).Return(nil)

	Index(render, repositories, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestIndexError(t *testing.T) {

	admin := &doorbot.Administrator{}
	accounts := []*doorbot.Account{}
	err := errors.New("i like pasta")

	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(1)

	repo.On("All", db).Return(accounts, err)

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Index(render, repositories, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestGet(t *testing.T) {

	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	account := &doorbot.Account{
		ID:   33,
		Name: "ACME",
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			ID:          3456,
			AccountType: doorbot.AccountManager,
		},
	}

	repo.On("Find", db, uint(33)).Return(account, nil)

	render.On("JSON", http.StatusOK, AccountViewModel{Account: account}).Return(nil)

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestGetNotOwner(t *testing.T) {

	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	account := &doorbot.Account{
		ID:   33,
		Name: "ACME",
	}

	resp := PublicAccount{
		ID:   33,
		Name: "ACME",
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationPerson,
		Person: &doorbot.Person{
			ID: 3456,
		},
	}

	render.On("JSON", http.StatusOK, PublicAccountViewModel{Account: resp}).Return(nil)
	repo.On("Find", db, uint(33)).Return(account, nil)

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestGetParseIntError(t *testing.T) {

	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)

	session := &auth.Authorization{}

	params := martini.Params{
		"id": "help",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestGetNotFound(t *testing.T) {

	var account *doorbot.Account

	session := &auth.Authorization{}

	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("DB").Return(db)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("AccountScope").Return(1)

	params := martini.Params{
		"id": "33",
	}

	repo.On("Find", db, uint(33)).Return(account, nil)

	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{})).Return()

	Get(render, repositories, params, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestPost(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)

	account := &doorbot.Account{
		Name: "ACME",
		Host: "derp",
	}

	admin := &doorbot.Administrator{}

	// nil
	var findByHostReponse *doorbot.Account

	repo.On("Create", db, account).Return(nil)
	repo.On("FindByHost", db, "derp").Return(findByHostReponse, nil)

	render.On("JSON", http.StatusCreated, AccountViewModel{Account: account}).Return()

	Post(render, repositories, AccountViewModel{Account: account}, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestPostCreateError(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	admin := &doorbot.Administrator{}

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(1)

	account := &doorbot.Account{
		Name: "ACME",
		Host: "derp",
	}

	// nil
	var findByHostReponse *doorbot.Account

	repo.On("Create", db, account).Return(errors.New("errooor"))
	repo.On("FindByHost", db, "derp").Return(findByHostReponse, nil)

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Post(render, repositories, AccountViewModel{Account: account}, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestPut(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)

	postAccount := &doorbot.Account{
		Name: "Romanian Landlords",
	}

	repoAccount := &doorbot.Account{
		ID:        5555,
		Name:      "ACME",
		IsEnabled: true,
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationAdministrator,
	}

	repo.On("Update", db, repoAccount).Return(true, nil)

	render.On("JSON", http.StatusOK, AccountViewModel{Account: repoAccount}).Return()

	Put(render, repoAccount, repositories, AccountViewModel{Account: postAccount}, session)

	assert.Equal(t, "Romanian Landlords", repoAccount.Name)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestPutFailed(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(1)

	postAccount := &doorbot.Account{
		Name: "Romanian Landlords",
	}

	repoAccount := &doorbot.Account{
		ID:        5555,
		Name:      "ACME",
		IsEnabled: true,
	}

	session := &auth.Authorization{
		Type: auth.AuthorizationAdministrator,
	}

	repo.On("Update", db, repoAccount).Return(false, errors.New("failed"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Put(render, repoAccount, repositories, AccountViewModel{Account: postAccount}, session)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	render := new(tests.MockRender)
	repo := new(tests.MockAccountRepository)
	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "33",
	}

	account := &doorbot.Account{
		ID:   33,
		Name: "ACME",
	}

	admin := &doorbot.Administrator{}

	repo.On("Find", db, uint(33)).Return(account, nil)
	repo.On("Delete", db, account).Return(true, nil)

	render.On("Status", http.StatusNoContent).Return()

	Delete(render, repositories, params, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestDeleteInvalidID(t *testing.T) {
	repo := new(tests.MockAccountRepository)
	render := new(tests.MockRender)

	admin := &doorbot.Administrator{}

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)

	params := martini.Params{
		"id": "help",
	}

	render.On("JSON", http.StatusBadRequest, doorbot.NewBadRequestErrorResponse([]string{"The id must be an unsigned integer"})).Return()

	Delete(render, repositories, params, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestDeleteNotFound(t *testing.T) {
	var account *doorbot.Account
	admin := &doorbot.Administrator{}

	repo := new(tests.MockAccountRepository)
	render := new(tests.MockRender)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("AccountScope").Return(1)

	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	params := martini.Params{
		"id": "44",
	}

	repo.On("Find", db, uint(44)).Return(account, nil)

	render.On("JSON", http.StatusNotFound, doorbot.NewEntityNotFoundResponse([]string{"The specified account does not exists."})).Return()

	Delete(render, repositories, params, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestDeleteFailed(t *testing.T) {
	repo := new(tests.MockAccountRepository)
	render := new(tests.MockRender)

	admin := &doorbot.Administrator{}

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(repo)
	repositories.On("DB").Return(db)
	repositories.On("AccountScope").Return(1)

	params := martini.Params{
		"id": "55",
	}

	account := &doorbot.Account{
		ID:   55,
		Name: "ACME",
	}

	repo.On("Find", db, uint(55)).Return(account, nil)
	repo.On("Delete", db, account).Return(false, errors.New("error"))

	render.On("JSON", http.StatusInternalServerError, doorbot.NewInternalServerErrorResponse([]string{})).Return()

	Delete(render, repositories, params, admin)

	render.Mock.AssertExpectations(t)
	repo.Mock.AssertExpectations(t)
}

func TestRegister(t *testing.T) {
	render := new(tests.MockRender)
	notificator := new(tests.MockNotificator)

	accountRepo := new(tests.MockAccountRepository)
	authRepo := new(tests.MockAuthenticationRepository)
	personRepo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)
	tx := new(tests.MockTransaction)

	config := &doorbot.DoorbotConfig{}

	repositories := new(tests.MockRepositories)
	repositories.On("AccountRepository").Return(accountRepo)
	repositories.On("PersonRepository").Return(personRepo)
	repositories.On("AuthenticationRepository").Return(authRepo)
	repositories.On("SetAccountScope", uint(0)).Return()

	repositories.On("DB").Return(db)
	repositories.On("Transaction").Return(tx, nil)

	vm := RegisterViewModel{
		Account: AccountRegisterRequest{
			Name: "ACME",
		},
	}

	var noAccount *doorbot.Account

	accountRepo.On("FindByHost", db, mock.AnythingOfType("string")).Return(noAccount, nil)
	accountRepo.On("Create", tx, mock.AnythingOfType("*doorbot.Account")).Return(nil)
	personRepo.On("Create", tx, mock.AnythingOfType("*doorbot.Person")).Return(nil)
	authRepo.On("Create", tx, mock.AnythingOfType("*doorbot.Authentication")).Return(nil)

	notificator.On("AccountCreated", mock.AnythingOfType("*doorbot.Account"), mock.AnythingOfType("*doorbot.Person"), mock.AnythingOfType("string")).Return()
	tx.On("Commit").Return(nil)

	render.On("JSON", http.StatusCreated, mock.AnythingOfType("AccountViewModel")).Return()

	Register(render, config, repositories, notificator, vm)

	render.Mock.AssertExpectations(t)
	accountRepo.Mock.AssertExpectations(t)
	personRepo.Mock.AssertExpectations(t)
	authRepo.Mock.AssertExpectations(t)
}
