package api

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/api/auth"
	"bitbucket.org/msamson/doorbot-api/tests"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthPassword_UserNotFound(t *testing.T) {

	requestBody, _ := json.Marshal(auth.PasswordRequest{
		Authentication: auth.PasswordAuthentication{
			Email:    "test@test.com",
			Password: "test",
		},
	})

	req, _ := http.NewRequest("POST", "/api/auth/password", bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "account.example.com")
	req.Header.Add("Authorization", "dashboard gatekeeper")

	rec := httptest.NewRecorder()
	server := newServer()

	account := &doorbot.Account{
		ID: 45,
	}

	var person *doorbot.Person

	accountRepo := new(tests.MockAccountRepository)
	authRepo := new(tests.MockAuthenticationRepository)
	personRepo := new(tests.MockPersonRepository)

	db := new(tests.MockExecutor)
	repos := getDependency(server, (*doorbot.Repositories)(nil)).(*tests.MockRepositories)
	repos.On("SetAccountScope", uint(45)).Return()
	repos.On("AccountRepository").Return(accountRepo)
	repos.On("AuthenticationRepository").Return(authRepo)
	repos.On("PersonRepository").Return(personRepo)
	repos.On("DB").Return(db)

	accountRepo.On("FindByHost", db, "account").Return(account, nil)
	personRepo.On("FindByEmail", db, "test@test.com").Return(person, nil)

	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	repos.AssertExpectations(t)
	personRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)

}
