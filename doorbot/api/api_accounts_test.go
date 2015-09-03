// +build tests

package api

import (
	"github.com/masom/doorbot/doorbot"
	"bitbucket.org/msamson/doorbot-api/tests"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountsIndex(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/accounts", nil)
	req.Header.Add("Host", "example.com")
	req.Header.Add("Authorization", "administrator 2")

	rec := httptest.NewRecorder()
	server := newServer()

	account := &doorbot.Account{
		ID: 3,
	}
	accounts := []*doorbot.Account{account}

	auth := &doorbot.AdministratorAuthentication{
		AdministratorID: 3,
	}

	administrator := &doorbot.Administrator{}

	db := new(tests.MockExecutor)

	accountRepo := new(tests.MockAccountRepository)
	adminRepo := new(tests.MockAdministratorRepository)
	adminAuthRepo := new(tests.MockAdministratorAuthenticationRepository)

	repositories := getDependency(server, (*doorbot.Repositories)(nil)).(*tests.MockRepositories)
	repositories.On("DB").Return(db)
	repositories.On("AccountRepository").Return(accountRepo)
	repositories.On("AdministratorRepository").Return(adminRepo)
	repositories.On("AdministratorAuthenticationRepository").Return(adminAuthRepo)

	accountRepo.On("All", db).Return(accounts, nil)

	adminRepo.On("Find", db, uint(3)).Return(administrator, nil)
	adminAuthRepo.On("FindByProviderIDAndToken", db, uint(2), "2").Return(auth, nil)

	server.ServeHTTP(rec, req)

	adminRepo.Mock.AssertExpectations(t)
	adminAuthRepo.Mock.AssertExpectations(t)
	accountRepo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
	// Should not be available.
	assert.Equal(t, http.StatusOK, rec.Code)
}
