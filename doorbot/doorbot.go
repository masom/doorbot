package doorbot

import (
	"database/sql"
)

const (
	// AccountMember access right
	AccountMember = 1
	// AccountManager access right
	AccountManager = 2
	// AccountOwner access right
	AccountOwner = 3

	// EventSignIn data point
	EventSignIn = 1
	// EventSignOut data point
	EventSignOut = 2

	// EventDoorAdded data point
	EventDoorAdded = 10
	// EventDoorRemoved data point
	EventDoorRemoved = 11
	// EventDoorUpdated data point
	EventDoorUpdated = 12

	// EventPersonAdded data point
	EventPersonAdded = 20
	// EventPersonRemoved data point
	EventPersonRemoved = 21
	// EventPersonUpdated data point
	EventPersonUpdated = 22

	// EventDeviceAdded data point
	EventDeviceAdded = 30
	// EventDeviceRemoved data point
	EventDeviceRemoved = 31
	// EventDeviceUpdated data point
	EventDeviceUpdated = 32
	// EventDeviceAssigned data point
	EventDeviceAssigned = 33
	// EventDeviceUnassigned data point
	EventDeviceUnassigned = 34
	// EventDeviceSignIn data point
	EventDeviceSignIn = 35
	// EventDeviceSignOut data point
	EventDeviceSignOut = 36

	// EventNotificationSent data point
	EventNotificationSent = 100
	// EventNotificationSMSSent data point
	EventNotificationSMSSent = 101
	// EventNotificationEmailSent data point
	EventNotificationEmailSent = 102
	// EventNotificationAppSent data point
	EventNotificationAppSent = 103
	// EventNotificationWebhookSent data point
	EventNotificationWebhookSent = 104
)

// Repositories interface
type Repositories interface {
	// Return the database instance
	DB() Executor

	// Return the account scope
	AccountScope() uint

	AccountRepository() AccountRepository
	AuthenticationRepository() AuthenticationRepository
	AdministratorRepository() AdministratorRepository
	AdministratorAuthenticationRepository() AdministratorAuthenticationRepository
	BridgeUserRepository() BridgeUserRepository
	DeviceRepository() DeviceRepository
	DoorRepository() DoorRepository
	EventRepository() EventRepository
	PersonRepository() PersonRepository

	SetAccountScope(uint)

	Transaction() (Transaction, error)
}

// AccountRepository repository interface
type AccountRepository interface {
	All(Executor) ([]*Account, error)
	Create(Executor, *Account) error
	Delete(Executor, *Account) (bool, error)
	Find(Executor, uint) (*Account, error)
	FindByHost(Executor, string) (*Account, error)
	Update(Executor, *Account) (bool, error)
}

// AuthenticationRepository repository interface
type AuthenticationRepository interface {
	Create(Executor, *Authentication) error
	Delete(Executor, *Authentication) (bool, error)
	FindByPersonID(Executor, uint) ([]*Authentication, error)
	FindByPersonIDAndProviderID(Executor, uint, uint) (*Authentication, error)
	FindByProviderIDAndPersonIDAndToken(Executor, uint, uint, string) (*Authentication, error)
	FindByProviderIDAndToken(Executor, uint, string) (*Authentication, error)
	Update(Executor, *Authentication) (bool, error)

	SetAccountScope(uint)
}

// AdministratorRepository repository interface
type AdministratorRepository interface {
	All(Executor) ([]*Administrator, error)
	Create(Executor, *Administrator) error
	Delete(Executor, *Administrator) (bool, error)
	Find(Executor, uint) (*Administrator, error)
	Update(Executor, *Administrator) (bool, error)
}

// AdministratorAuthenticationRepository repository interface
type AdministratorAuthenticationRepository interface {
	All(Executor) ([]*AdministratorAuthentication, error)
	Create(Executor, *AdministratorAuthentication) error
	Delete(Executor, *AdministratorAuthentication) (bool, error)
	FindByAdministratorID(Executor, uint) ([]*AdministratorAuthentication, error)
	FindByAdministratorIDAndProviderID(e Executor, administrator uint, provider uint) (*AdministratorAuthentication, error)
	FindByProviderIDAndToken(e Executor, provider uint, token string) (*AdministratorAuthentication, error)
	Update(Executor, *AdministratorAuthentication) (bool, error)
}

// BridgeUserRepository repository interface
type BridgeUserRepository interface {
	FindByPersonID(Executor, uint) ([]*BridgeUser, error)
	FindByPersonIDAndBridgeID(t Executor, pID uint, bID uint) (*BridgeUser, error)
	FindByBridgeID(Executor, uint) ([]*BridgeUser, error)
	Create(Executor, *BridgeUser) error
	Delete(Executor, *BridgeUser) (bool, error)
	SetAccountScope(Executor uint)
}

// DeviceRepository repository interface
type DeviceRepository interface {
	All(Executor) ([]*Device, error)
	Create(Executor, *Device) error
	Delete(Executor, *Device) (bool, error)
	Find(Executor, uint) (*Device, error)
	FindByDeviceID(Executor, string) (*Device, error)
	FindByToken(Executor, string) (*Device, error)
	SetAccountScope(uint)
	Update(t Executor, d *Device) (bool, error)
	Enable(t Executor, d *Device, enabled bool) (bool, error)
}

// DoorRepository repository interface
type DoorRepository interface {
	All(Executor) ([]*Door, error)
	Create(Executor, *Door) error
	Delete(Executor, *Door) (bool, error)
	Find(Executor, uint) (*Door, error)
	SetAccountScope(uint)
	Update(Executor, *Door) (bool, error)
}

// EventRepository repository interface
type EventRepository interface {
	All(Executor, *Pagination) ([]*Event, error)
	Create(Executor, *Event) error
	Find(Executor, uint) (*Event, error)
	SetAccountScope(uint)
}

// PersonRepository repository interface
type PersonRepository interface {
	All(Executor) ([]*Person, error)
	Create(Executor, *Person) error
	Delete(Executor, *Person) (bool, error)
	Find(Executor, uint) (*Person, error)
	FindByEmail(Executor, string) (*Person, error)
	SetAccountScope(uint)
	Update(Executor, *Person) (bool, error)
}

// Pagination data structure
type Pagination struct {
	Page  uint
	Limit uint
	Order string
}

// Notificator interface
type Notificator interface {
	Notify(a *Account, d *Door, p *Person) error
	AccountCreated(a *Account, p *Person, password string)
}

// Transaction interface for database
type Transaction interface {
	Commit() error
	Rollback() error

	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectInt(query string, args ...interface{}) (int64, error)
	SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)
	SelectFloat(query string, args ...interface{}) (float64, error)
	SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error)
	SelectStr(query string, args ...interface{}) (string, error)
	SelectNullStr(query string, args ...interface{}) (sql.NullString, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
}

// Executor database interface
type Executor interface {
	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectInt(query string, args ...interface{}) (int64, error)
	SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)
	SelectFloat(query string, args ...interface{}) (float64, error)
	SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error)
	SelectStr(query string, args ...interface{}) (string, error)
	SelectNullStr(query string, args ...interface{}) (sql.NullString, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
}
