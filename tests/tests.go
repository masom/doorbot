package tests

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"database/sql"
	"github.com/stretchr/testify/mock"
	"github.com/martini-contrib/render"
	"github.com/masom/doorbot/doorbot"
)

type MockNotificator struct {
	mock.Mock
}

func (m *MockNotificator) AccountCreated(a *doorbot.Account, p *doorbot.Person, password string) {
	m.Mock.Called(a, p, password)
}

func (m *MockNotificator) Notify(a *doorbot.Account, d *doorbot.Door, p *doorbot.Person) error {
	args := m.Mock.Called(a, d, p)
	return args.Error(0)
}

type MockBridges struct {
	mock.Mock
}

func (m *MockBridges) GetUsers(bid uint) ([]*doorbot.BridgeUser, error) {
	args := m.Mock.Called(bid)
	return args.Get(0).([]*doorbot.BridgeUser), args.Error(1)
}

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	a := m.Mock.Called(i, keys)
	return a.Get(0).(interface{}), a.Error(1)
}

func (m *MockExecutor) Insert(list ...interface{}) error {
	a := m.Mock.Called(list)
	return a.Error(0)
}

func (m *MockExecutor) Update(list ...interface{}) (int64, error) {
	a := m.Mock.Called(list)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockExecutor) Delete(list ...interface{}) (int64, error) {
	a := m.Mock.Called(list)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.Result), a.Error(1)
}

func (m *MockExecutor) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error){
	a := m.Mock.Called(i, query, args)
	return a.Get(0).([]interface{}), a.Error(1)

}

func (m *MockExecutor) SelectInt(query string, args ...interface{}) (int64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockExecutor) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.NullInt64), a.Error(1)
}

func (m *MockExecutor) SelectFloat(query string, args ...interface{}) (float64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(float64), a.Error(1)
}

func (m *MockExecutor) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.NullFloat64), a.Error(1)
}
func (m *MockExecutor) SelectStr(query string, args ...interface{}) (string, error) {
	a := m.Mock.Called(query, args)
	return a.String(0), a.Error(1)
}

func (m *MockExecutor) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.NullString), a.Error(1)
}

func (m *MockExecutor) SelectOne(holder interface{}, query string, args ...interface{}) error {
	a := m.Mock.Called(holder, query, args)
	return a.Error(0)
}

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Begin() error {
	return m.Mock.Called().Error(0)
}

func (m *MockTransaction) Commit() error {
	return m.Mock.Called().Error(0)
}

func (m *MockTransaction) Rollback() error {
	return m.Mock.Called().Error(0)
}

func (m *MockTransaction) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	a := m.Mock.Called(i, keys)
	return a.Get(0).(interface{}), a.Error(1)
}

func (m *MockTransaction) Insert(list ...interface{}) error {
	a := m.Mock.Called(list)
	return a.Error(0)
}

func (m *MockTransaction) Update(list ...interface{}) (int64, error) {
	a := m.Mock.Called(list)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockTransaction) Delete(list ...interface{}) (int64, error) {
	a := m.Mock.Called(list)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.Result), a.Error(1)
}

func (m *MockTransaction) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error){
	a := m.Mock.Called(i, query, args)
	return a.Get(0).([]interface{}), a.Error(1)

}

func (m *MockTransaction) SelectInt(query string, args ...interface{}) (int64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockTransaction) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.NullInt64), a.Error(1)
}

func (m *MockTransaction) SelectFloat(query string, args ...interface{}) (float64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(float64), a.Error(1)
}

func (m *MockTransaction) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.NullFloat64), a.Error(1)
}
func (m *MockTransaction) SelectStr(query string, args ...interface{}) (string, error) {
	a := m.Mock.Called(query, args)
	return a.String(0), a.Error(1)
}

func (m *MockTransaction) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	a := m.Mock.Called(query, args)
	return a.Get(0).(sql.NullString), a.Error(1)
}

func (m *MockTransaction) SelectOne(holder interface{}, query string, args ...interface{}) error {
	a := m.Mock.Called(holder, query, args)
	return a.Error(0)
}

type MockRepositories struct {
	mock.Mock
}

func (m *MockRepositories) AccountScope() uint {
	args := m.Mock.Called()
	return args.Get(0).(uint)
}

func (m *MockRepositories) Transaction() (doorbot.Transaction, error) {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.Transaction), args.Error(1)
}

func (m *MockRepositories) AccountRepository() doorbot.AccountRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.AccountRepository)
}

func (m *MockRepositories) AdministratorRepository() doorbot.AdministratorRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.AdministratorRepository)
}

func (m *MockRepositories) AdministratorAuthenticationRepository() doorbot.AdministratorAuthenticationRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.AdministratorAuthenticationRepository)
}

func (m *MockRepositories) AuthenticationRepository() doorbot.AuthenticationRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.AuthenticationRepository)
}

func (m *MockRepositories) BridgeUserRepository() doorbot.BridgeUserRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.BridgeUserRepository)
}

func (m *MockRepositories) DeviceRepository() doorbot.DeviceRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.DeviceRepository)
}

func (m *MockRepositories) DB() doorbot.Executor {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.Executor)
}

func (m *MockRepositories) SetAccountScope(id uint) {
	m.Mock.Called(id)
}

func (m *MockRepositories) DoorRepository() doorbot.DoorRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.DoorRepository)
}

func (m *MockRepositories) EventRepository() doorbot.EventRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.EventRepository)
}

func (m *MockRepositories) PersonRepository() doorbot.PersonRepository {
	args := m.Mock.Called()
	return args.Get(0).(doorbot.PersonRepository)
}



type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) All(e doorbot.Executor) ([]*doorbot.Account, error) {
	args := m.Mock.Called(e)

	return args.Get(0).([]*doorbot.Account), args.Error(1)
}

func (m *MockAccountRepository) Find(e doorbot.Executor, id uint) (*doorbot.Account, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).(*doorbot.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByHost(e doorbot.Executor, host string) (*doorbot.Account, error) {
	args := m.Mock.Called(e, host)

	return args.Get(0).(*doorbot.Account), args.Error(1)
}

func (m *MockAccountRepository) Create(e doorbot.Executor, a *doorbot.Account) error {
	return m.Mock.Called(e, a).Error(0)
}

func (m *MockAccountRepository) Update(e doorbot.Executor, a *doorbot.Account) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockAccountRepository) Delete(e doorbot.Executor, a *doorbot.Account) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

type MockAuthenticationRepository struct {
	mock.Mock
}

func (m *MockAuthenticationRepository) FindByPersonID(e doorbot.Executor, id uint) ([]*doorbot.Authentication, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).([]*doorbot.Authentication), args.Error(1)
}

func (m *MockAuthenticationRepository) FindByPersonIDAndProviderID(e doorbot.Executor, personID uint, providerID uint) (*doorbot.Authentication, error) {
	args := m.Mock.Called(e, personID, providerID)
	return args.Get(0).(*doorbot.Authentication), args.Error(1)
}

func (m *MockAuthenticationRepository) FindByProviderIDAndToken(e doorbot.Executor, providerID uint, token string) (*doorbot.Authentication, error) {
	args := m.Mock.Called(e, providerID, token)
	return args.Get(0).(*doorbot.Authentication), args.Error(1)
}

func (m *MockAuthenticationRepository) FindByProviderIDAndPersonIDAndToken(e doorbot.Executor, providerID uint, personID uint, token string) (*doorbot.Authentication, error) {
	args := m.Mock.Called(e, providerID, personID, token)
	return args.Get(0).(*doorbot.Authentication), args.Error(1)
}

func (m *MockAuthenticationRepository) Create(e doorbot.Executor, a *doorbot.Authentication) error {
	return m.Mock.Called(e, a).Error(0)
}

func (m *MockAuthenticationRepository) Update(e doorbot.Executor, a *doorbot.Authentication) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthenticationRepository) Delete(e doorbot.Executor, a *doorbot.Authentication) (bool, error) {
	args := m.Mock.Called(e,a)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthenticationRepository) SetAccountScope(accountID uint) {
	m.Mock.Called(accountID)
}

type MockAdministratorRepository struct {
	mock.Mock
}

func (m *MockAdministratorRepository) All(e doorbot.Executor) ([]*doorbot.Administrator, error) {
	args := m.Mock.Called(e)
	return args.Get(0).([]*doorbot.Administrator), args.Error(1)
}

func (m *MockAdministratorRepository) Create(e doorbot.Executor, a *doorbot.Administrator) error {
	args := m.Mock.Called(e)
	return args.Error(0)
}

func (m *MockAdministratorRepository) Delete(e doorbot.Executor, a *doorbot.Administrator) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockAdministratorRepository) Find(e doorbot.Executor, id uint) (*doorbot.Administrator, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).(*doorbot.Administrator), args.Error(1)
}

func (m *MockAdministratorRepository) Update(e doorbot.Executor, a *doorbot.Administrator) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

type MockAdministratorAuthenticationRepository struct {
	mock.Mock
}

func (m *MockAdministratorAuthenticationRepository) All(e doorbot.Executor) ([]*doorbot.AdministratorAuthentication, error) {
	args := m.Mock.Called(e)
	return args.Get(0).([]*doorbot.AdministratorAuthentication), args.Error(1)
}

func (m *MockAdministratorAuthenticationRepository) Create(e doorbot.Executor, a *doorbot.AdministratorAuthentication) error {
	args := m.Mock.Called(e)
	return args.Error(0)
}

func (m *MockAdministratorAuthenticationRepository) Delete(e doorbot.Executor, a *doorbot.AdministratorAuthentication) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockAdministratorAuthenticationRepository) Update(e doorbot.Executor, a *doorbot.AdministratorAuthentication) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockAdministratorAuthenticationRepository) FindByAdministratorID(e doorbot.Executor, id uint) ([]*doorbot.AdministratorAuthentication, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).([]*doorbot.AdministratorAuthentication), args.Error(1)
}

func (m *MockAdministratorAuthenticationRepository) FindByAdministratorIDAndProviderID(e doorbot.Executor, aid uint, provider uint) (*doorbot.AdministratorAuthentication, error) {
	args := m.Mock.Called(e, aid, provider)
	return args.Get(0).(*doorbot.AdministratorAuthentication), args.Error(1)
}

func (m *MockAdministratorAuthenticationRepository) FindByProviderIDAndToken(e doorbot.Executor, pid uint, token string) (*doorbot.AdministratorAuthentication, error) {
	args := m.Mock.Called(e, pid, token)
	return args.Get(0).(*doorbot.AdministratorAuthentication), args.Error(1)
}

type MockBridgeUserRepository struct {
	mock.Mock
}

func (m *MockBridgeUserRepository) SetAccountScope(s uint) {
	m.Mock.Called(s)
}

func (m *MockBridgeUserRepository) FindByPersonIDAndBridgeID(e doorbot.Executor, pid uint, bid uint) (*doorbot.BridgeUser, error) {
	args := m.Mock.Called(e, pid, bid)

	return args.Get(0).(*doorbot.BridgeUser), args.Error(1)
}

func (m *MockBridgeUserRepository) FindByPersonID(e doorbot.Executor, pid uint) ([]*doorbot.BridgeUser, error) {
	args := m.Mock.Called(e, pid)

	return args.Get(0).([]*doorbot.BridgeUser), args.Error(1)
}

func (m *MockBridgeUserRepository) FindByBridgeID(e doorbot.Executor, bid uint) ([]*doorbot.BridgeUser, error) {
	args := m.Mock.Called(e, bid)
	return args.Get(0).([]*doorbot.BridgeUser), args.Error(1)
}

func (m *MockBridgeUserRepository) Create(e doorbot.Executor, a *doorbot.BridgeUser) error {
	return m.Mock.Called(e, a).Error(0)
}

func (m *MockBridgeUserRepository) Update(e doorbot.Executor, a *doorbot.BridgeUser) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockBridgeUserRepository) Delete(e doorbot.Executor, a *doorbot.BridgeUser) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) All(e doorbot.Executor) ([]*doorbot.Device, error) {
	args := m.Mock.Called(e)

	return args.Get(0).([]*doorbot.Device), args.Error(1)
}

func (m *MockDeviceRepository) Find(e doorbot.Executor, id uint) (*doorbot.Device, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).(*doorbot.Device), args.Error(1)
}

func (m *MockDeviceRepository) FindByDeviceID(e doorbot.Executor, deviceID string) (*doorbot.Device, error) {
	args := m.Mock.Called(e, deviceID)

	return args.Get(0).(*doorbot.Device), args.Error(1)
}

func (m *MockDeviceRepository) FindByToken(e doorbot.Executor, token string) (*doorbot.Device, error) {
	args := m.Mock.Called(e, token)

	return args.Get(0).(*doorbot.Device), args.Error(1)
}

func (m *MockDeviceRepository) Create(e doorbot.Executor, a *doorbot.Device) error {
	return m.Mock.Called(e, a).Error(0)
}

func (m *MockDeviceRepository) Update(e doorbot.Executor, a *doorbot.Device) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockDeviceRepository) Delete(e doorbot.Executor, a *doorbot.Device) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockDeviceRepository) Enable(e doorbot.Executor, d *doorbot.Device, enabled bool) (bool, error) {
	args := m.Mock.Called(e, d, enabled)
	return args.Bool(0), args.Error(1)
}

func (m *MockDeviceRepository) SetAccountScope(accountID uint) {
	m.Mock.Called(accountID)
}

type MockDoorRepository struct {
	mock.Mock
}

func (m *MockDoorRepository) All(e doorbot.Executor) ([]*doorbot.Door, error) {
	args := m.Mock.Called(e)

	return args.Get(0).([]*doorbot.Door), args.Error(1)
}

func (m *MockDoorRepository) Find(e doorbot.Executor, id uint) (*doorbot.Door, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).(*doorbot.Door), args.Error(1)
}

func (m *MockDoorRepository) Create(e doorbot.Executor, a *doorbot.Door) error {
	return m.Mock.Called(e, a).Error(0)
}

func (m *MockDoorRepository) Update(e doorbot.Executor, a *doorbot.Door) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockDoorRepository) Delete(e doorbot.Executor, a *doorbot.Door) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockDoorRepository) SetAccountScope(accountID uint) {
	m.Mock.Called(accountID)
}

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) All(e doorbot.Executor, p *doorbot.Pagination) ([]*doorbot.Event, error) {
	args := m.Mock.Called(e, p)
	return args.Get(0).([]*doorbot.Event), args.Error(1)
}

func (m *MockEventRepository) Create(ex doorbot.Executor, e *doorbot.Event) error {
	args := m.Mock.Called(ex, e)
	return args.Error(0)
}

func (m *MockEventRepository) Find(e doorbot.Executor, id uint) (*doorbot.Event, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).(*doorbot.Event), args.Error(1)
}

func (m *MockEventRepository) SetAccountScope(id uint) {
	m.Mock.Called(id)
}

type MockPersonRepository struct {
	mock.Mock
}

func (m *MockPersonRepository) All(e doorbot.Executor) ([]*doorbot.Person, error) {
	args := m.Mock.Called(e)

	return args.Get(0).([]*doorbot.Person), args.Error(1)
}

func (m *MockPersonRepository) Find(e doorbot.Executor, id uint) (*doorbot.Person, error) {
	args := m.Mock.Called(e, id)
	return args.Get(0).(*doorbot.Person), args.Error(1)
}

func (m *MockPersonRepository) FindByEmail(e doorbot.Executor, email string) (*doorbot.Person, error) {
	args := m.Mock.Called(e, email)

	return args.Get(0).(*doorbot.Person), args.Error(1)
}

func (m *MockPersonRepository) Create(e doorbot.Executor, a *doorbot.Person) error {
	return m.Mock.Called(e, a).Error(0)
}

func (m *MockPersonRepository) Update(e doorbot.Executor, a *doorbot.Person) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockPersonRepository) Delete(e doorbot.Executor, a *doorbot.Person) (bool, error) {
	args := m.Mock.Called(e, a)
	return args.Bool(0), args.Error(1)
}

func (m *MockPersonRepository) SetAccountScope(accountID uint) {
	m.Mock.Called(accountID)
}

type MockRender struct {
	mock.Mock
}

func (m *MockRender) JSON(status int, v interface{}) {
	m.Called(status, v)
}

func (m *MockRender) HTML(status int, name string, v interface{}, htmlOpt ...render.HTMLOptions) {
	m.Called(status, name, v, htmlOpt)
}

func (m *MockRender) XML(status int, v interface{}) {
	m.Called(status, v)
}

func (m *MockRender) Data(status int, v []byte) {
	m.Called(status, v)
}

func (m *MockRender) Error(status int) {
	m.Called(status)
}

func (m *MockRender) Status(status int) {
	m.Called(status)
}

func (m *MockRender) Redirect(location string, status ...int) {
	m.Called(location, status)
}

func (m *MockRender) Template() *template.Template {
	m.Called()
	return &template.Template{}
}

func (m *MockRender) Header() http.Header {

	r := httptest.ResponseRecorder{}
	m.Called();

	return r.Header()
}
