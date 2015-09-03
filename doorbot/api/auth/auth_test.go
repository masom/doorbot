package auth

import (
	"github.com/masom/doorbot/doorbot"
	"github.com/masom/doorbot/doorbot/auth"
	"bitbucket.org/msamson/doorbot-api/tests"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestPassword(t *testing.T) {
	account := &doorbot.Account{
		ID:        444,
		Name:      "ACME",
		IsEnabled: true,
	}

	person := &doorbot.Person{
		ID:    1,
		Name:  "Cookie Monster",
		Email: "cookiemonster@aol.com",
	}

	passwordAuthentication := &doorbot.Authentication{
		AccountID: 444,
		PersonID:  1,
		Token:     "$2a$10$8XdprxFRIXCv1TC2cDjMNuQRiYkOX9PIivVpnSMM9b.1UjulLlrVm", // test
	}

	passwordRequest := PasswordRequest{
		Authentication: PasswordAuthentication{
			Email:    "cookiemonster@aol.com",
			Password: "test",
		},
	}

	var tokenNotFound *doorbot.Authentication

	render := new(tests.MockRender)
	personRepo := new(tests.MockPersonRepository)
	authRepo := new(tests.MockAuthenticationRepository)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(personRepo)
	repositories.On("AuthenticationRepository").Return(authRepo)

	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	personRepo.On("FindByEmail", db, "cookiemonster@aol.com").Return(person, nil)

	authRepo.On("FindByPersonIDAndProviderID", db, person.ID, auth.ProviderPassword).Return(passwordAuthentication, nil).Once()
	authRepo.On("FindByPersonIDAndProviderID", db, person.ID, auth.ProviderAPIToken).Return(tokenNotFound, nil).Once()
	authRepo.On("Create", db, mock.AnythingOfType("*doorbot.Authentication")).Return(nil).Once()

	render.On("JSON", http.StatusOK, mock.AnythingOfType("auth.APITokenResponse")).Return().Once()
	Password(render, account, repositories, passwordRequest)

	render.Mock.AssertExpectations(t)
	personRepo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
	authRepo.Mock.AssertExpectations(t)
}

func TestPasswordReuseToken(t *testing.T) {

	account := &doorbot.Account{
		ID:        1,
		Name:      "ACME",
		IsEnabled: true,
	}

	person := &doorbot.Person{
		ID:    1,
		Name:  "Cookie Monster",
		Email: "cookiemonster@aol.com",
	}

	passwordAuthentication := &doorbot.Authentication{
		AccountID: 1,
		PersonID:  1,
		Token:     "$2a$10$8XdprxFRIXCv1TC2cDjMNuQRiYkOX9PIivVpnSMM9b.1UjulLlrVm", // test
	}

	apiTokenAuthentication := &doorbot.Authentication{
		AccountID: 1,
		PersonID:  1,
		Token:     "i-like-pasta",
	}

	passwordRequest := PasswordRequest{
		Authentication: PasswordAuthentication{
			Email:    "cookiemonster@aol.com",
			Password: "test",
		},
	}

	token := APITokenResponse{
		Authentication: TokenAuthentication{Token: "1.i-like-pasta"},
		Person:         person,
	}

	render := new(tests.MockRender)
	personRepo := new(tests.MockPersonRepository)
	authRepo := new(tests.MockAuthenticationRepository)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(personRepo)
	repositories.On("AuthenticationRepository").Return(authRepo)

	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	personRepo.On("FindByEmail", db, "cookiemonster@aol.com").Return(person, nil)
	authRepo.On("FindByPersonIDAndProviderID", db, person.ID, auth.ProviderPassword).Return(passwordAuthentication, nil)
	authRepo.On("FindByPersonIDAndProviderID", db, person.ID, auth.ProviderAPIToken).Return(apiTokenAuthentication, nil)

	render.On("JSON", http.StatusOK, token).Return()

	Password(render, account, repositories, passwordRequest)

	render.Mock.AssertExpectations(t)
	personRepo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
	authRepo.Mock.AssertExpectations(t)
}

func TestPasswordPersonNotFound(t *testing.T) {
	var person *doorbot.Person

	account := &doorbot.Account{
		ID:        1,
		Name:      "ACME",
		IsEnabled: true,
	}

	passwordRequest := PasswordRequest{
		Authentication: PasswordAuthentication{
			Email:    "cookiemonster@aol.com",
			Password: "test",
		},
	}

	render := new(tests.MockRender)
	personRepo := new(tests.MockPersonRepository)
	authRepo := new(tests.MockAuthenticationRepository)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(personRepo)
	repositories.On("AuthenticationRepository").Return(authRepo)

	db := new(tests.MockExecutor)
	repositories.On("DB").Return(db)

	personRepo.On("FindByEmail", db, "cookiemonster@aol.com").Return(person, nil)
	render.On("JSON", http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"})).Return()
	Password(render, account, repositories, passwordRequest)

	render.Mock.AssertExpectations(t)
	personRepo.Mock.AssertExpectations(t)
	authRepo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}

func TestPasswordInvalid(t *testing.T) {
	person := &doorbot.Person{
		ID:    1,
		Name:  "Cookie Monster",
		Email: "cookiemonster@aol.com",
	}

	account := &doorbot.Account{
		ID:        1,
		Name:      "ACME",
		IsEnabled: true,
	}

	passwordAuthentication := &doorbot.Authentication{
		AccountID: 1,
		PersonID:  1,
		Token:     "$2a$10$8XdprxFRIXCv1TC2cDjMNuQRiYkOX9PIivVpnSMM9b.1UjulLlrVm", // test
	}

	passwordRequest := PasswordRequest{
		Authentication: PasswordAuthentication{
			Email:    "cookiemonster@aol.com",
			Password: "test1",
		},
	}

	render := new(tests.MockRender)
	personRepo := new(tests.MockPersonRepository)
	authRepo := new(tests.MockAuthenticationRepository)

	db := new(tests.MockExecutor)

	repositories := new(tests.MockRepositories)
	repositories.On("PersonRepository").Return(personRepo)
	repositories.On("AuthenticationRepository").Return(authRepo)
	repositories.On("DB").Return(db)

	personRepo.On("FindByEmail", db, "cookiemonster@aol.com").Return(person, nil)
	authRepo.On("FindByPersonIDAndProviderID", db, person.ID, auth.ProviderPassword).Return(passwordAuthentication, nil)

	render.On("JSON", http.StatusUnauthorized, doorbot.NewUnauthorizedErrorResponse([]string{"Invalid email or password"})).Return()
	Password(render, account, repositories, passwordRequest)

	render.Mock.AssertExpectations(t)
	personRepo.Mock.AssertExpectations(t)
	authRepo.Mock.AssertExpectations(t)
	repositories.Mock.AssertExpectations(t)
}
