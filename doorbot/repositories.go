package doorbot

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"gopkg.in/gorp.v1"
	// Ensure lib/pq is not skipped by Godeps
	_ "github.com/lib/pq"
	"strconv"
)

type repositories struct {
	AccountID uint
	db        *gorp.DbMap

	account                     AccountRepository
	administrator               AdministratorRepository
	administratorAuthentication AdministratorAuthenticationRepository
	authentication              AuthenticationRepository
	bridgeUser                  BridgeUserRepository
	door                        DoorRepository
	device                      DeviceRepository
	event                       EventRepository
	person                      PersonRepository
}

// DB returns the database instance
func (r *repositories) DB() Executor {
	return r.db
}

// AccountScope returns the account scope
func (r *repositories) AccountScope() uint {
	return r.AccountID
}

// SetAccountScope sets the account scope to the provided value on all repositories
func (r *repositories) SetAccountScope(a uint) {
	r.AccountID = a

	if r.authentication != nil {
		r.authentication.SetAccountScope(a)
	}

	if r.bridgeUser != nil {
		r.bridgeUser.SetAccountScope(a)
	}

	if r.device != nil {
		r.device.SetAccountScope(a)
	}

	if r.door != nil {
		r.door.SetAccountScope(a)
	}

	if r.event != nil {
		r.event.SetAccountScope(a)
	}

	if r.person != nil {
		r.person.SetAccountScope(a)
	}
}

// Transaction creates a new Transaction
func (r *repositories) Transaction() (Transaction, error) {
	return r.db.Begin()
}

type accountRepository struct {
	AccountID uint
}

type administratorRepository struct {
}

type administratorAuthenticationRepository struct {
}

type authenticationRepository struct {
	AccountID uint
}

type bridgeUserRepository struct {
	AccountID uint
}

type doorRepository struct {
	AccountID uint
}

type deviceRepository struct {
	AccountID uint
}

type eventRepository struct {
	AccountID uint
}

type personRepository struct {
	AccountID uint
}

func newRepositories(d *gorp.DbMap) *repositories {
	return &repositories{
		db: d,
	}
}

// AccountRepository returns an AccountRepository instance
func (r *repositories) AccountRepository() AccountRepository {
	if r.account == nil {
		r.account = &accountRepository{}
	}

	return r.account
}

// AdministratorRepository returns an AdministratorRepository instance
func (r *repositories) AdministratorRepository() AdministratorRepository {
	if r.administrator == nil {
		r.administrator = &administratorRepository{}
	}

	return r.administrator
}

// AdministratorAuthenticationRepository returns an AdministratorAuthenticationRepository instance
func (r *repositories) AdministratorAuthenticationRepository() AdministratorAuthenticationRepository {
	if r.administratorAuthentication == nil {
		r.administratorAuthentication = &administratorAuthenticationRepository{}
	}
	return r.administratorAuthentication
}

// AuthenticationRepository returns an AuthenticationRepository instance
func (r *repositories) AuthenticationRepository() AuthenticationRepository {
	if r.authentication == nil {
		r.authentication = &authenticationRepository{
			AccountID: r.AccountID,
		}
	}
	return r.authentication
}

// BridgeUserRepository returns a BridgeUserRepository instance
func (r *repositories) BridgeUserRepository() BridgeUserRepository {
	if r.bridgeUser == nil {
		r.bridgeUser = &bridgeUserRepository{
			AccountID: r.AccountID,
		}
	}

	return r.bridgeUser
}

// DoorRepository returns a DoorRepository instance
func (r *repositories) DoorRepository() DoorRepository {
	if r.door == nil {
		r.door = &doorRepository{
			AccountID: r.AccountID,
		}
	}
	return r.door
}

// PersonRepository returns a PersonRepository instance
func (r *repositories) PersonRepository() PersonRepository {
	if r.person == nil {
		r.person = &personRepository{
			AccountID: r.AccountID,
		}
	}

	return r.person
}

// DeviceRepository returns a DeviceRepository instance
func (r *repositories) DeviceRepository() DeviceRepository {
	if r.device == nil {
		r.device = &deviceRepository{
			AccountID: r.AccountID,
		}
	}
	return r.device
}

// EventRepository retuns a EventRepository instance
func (r *repositories) EventRepository() EventRepository {
	if r.event == nil {
		r.event = &eventRepository{
			AccountID: r.AccountID,
		}
	}
	return r.event
}

// All returns all accounts
func (r *accountRepository) All(t Executor) ([]*Account, error) {
	var accounts []*Account

	_, err := t.Select(
		&accounts,
		"SELECT * FROM accounts ORDER BY name ASC",
	)

	return accounts, err
}

// Find a specific account by id
func (r *accountRepository) Find(t Executor, id uint) (*Account, error) {
	var (
		accounts []*Account
		account  *Account
	)

	_, err := t.Select(
		&accounts,
		"SELECT * FROM accounts WHERE id = :id",
		map[string]interface{}{"id": id},
	)

	if len(accounts) == 1 {
		account = accounts[0]
	}

	return account, err
}

// FindByHost search for an account registered with the specified host.
// If the host is an integer the account is looked up by id.
func (r *accountRepository) FindByHost(t Executor, host string) (*Account, error) {
	var account *Account
	var accounts []*Account

	query := "SELECT * FROM accounts WHERE "
	parameters := map[string]interface{}{}

	_, err := strconv.ParseUint(host, 10, 32)
	if err == nil {
		query += "id = :id"
		parameters["id"] = host
	} else {
		query += "host = :host"
		parameters["host"] = host
	}

	query += " LIMIT 1"

	_, err = t.Select(
		&accounts,
		query,
		parameters,
	)

	if len(accounts) == 1 {
		account = accounts[0]
	}

	return account, err
}

func (r *accountRepository) Create(t Executor, account *Account) error {
	return t.Insert(account)
}

func (r *accountRepository) Update(t Executor, account *Account) (bool, error) {
	count, err := t.Update(account)
	return count > 0, err
}

func (r *accountRepository) Delete(t Executor, account *Account) (bool, error) {
	count, err := t.Delete(account)
	return count > 0, err
}

func (r *accountRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func (r *administratorRepository) All(t Executor) ([]*Administrator, error) {
	var administrators []*Administrator

	_, err := t.Select(
		&administrators,
		"SELECT * FROM administrators ORDER BY NAME ASC",
	)

	return administrators, err
}

func (r *administratorRepository) Find(t Executor, id uint) (*Administrator, error) {
	var (
		administrators []*Administrator
		administrator  *Administrator
	)

	_, err := t.Select(
		&administrators,
		"SELECT * FROM administrators WHERE id = :id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"id": id},
	)

	if len(administrators) == 1 {
		administrator = administrators[0]
	}

	return administrator, err
}

func (r *administratorRepository) Create(t Executor, administrator *Administrator) error {
	return t.Insert(administrator)
}

func (r *administratorRepository) Update(t Executor, administrator *Administrator) (bool, error) {
	count, err := t.Update(administrator)
	return count > 0, err
}

func (r *administratorRepository) Delete(t Executor, administrator *Administrator) (bool, error) {
	count, err := t.Delete(administrator)
	return count > 0, err
}

func (r *administratorAuthenticationRepository) All(t Executor) ([]*AdministratorAuthentication, error) {
	var administratorAuthentications []*AdministratorAuthentication

	_, err := t.Select(
		&administratorAuthentications,
		"SELECT * FROM administratorAuthentications ORDER BY NAME ASC",
	)

	return administratorAuthentications, err
}

func (r *administratorAuthenticationRepository) FindByAdministratorID(t Executor, administratorID uint) ([]*AdministratorAuthentication, error) {
	var administratorAuthentications []*AdministratorAuthentication

	_, err := t.Select(
		&administratorAuthentications,
		"SELECT * FROM administratorAuthentications WHERE person_id = :id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"administrator_id": administratorID},
	)

	return administratorAuthentications, err
}

func (r *administratorAuthenticationRepository) FindByAdministratorIDAndProviderID(t Executor, administratorID uint, providerID uint) (*AdministratorAuthentication, error) {
	var (
		administratorAuthentications []*AdministratorAuthentication
		administratorAuthentication  *AdministratorAuthentication
	)

	_, err := t.Select(
		&administratorAuthentications,
		"SELECT * FROM administratorAuthentications WHERE administrator_id = :administrator_id AND provider_id = :provider_id LIMIT 1",
		map[string]interface{}{
			"administrator_id": administratorID,
			"provider_id":      providerID,
		},
	)

	if len(administratorAuthentications) == 1 {
		administratorAuthentication = administratorAuthentications[0]
	}

	return administratorAuthentication, err
}

func (r *administratorAuthenticationRepository) FindByProviderIDAndToken(e Executor, pid uint, token string) (*AdministratorAuthentication, error) {
	var (
		administratorAuthentications []*AdministratorAuthentication
		administratorAuthentication  *AdministratorAuthentication
	)

	_, err := e.Select(
		&administratorAuthentications,
		"SELECT * FROM administratorAuthentications WHERE token = :token AND provider_id = :provider_id LIMIT 1",
		map[string]interface{}{
			"token":       token,
			"provider_id": pid,
		},
	)

	if len(administratorAuthentications) == 1 {
		administratorAuthentication = administratorAuthentications[0]
	}

	return administratorAuthentication, err
}

func (r *administratorAuthenticationRepository) Create(t Executor, administratorAuthentication *AdministratorAuthentication) error {
	return t.Insert(administratorAuthentication)
}

func (r *administratorAuthenticationRepository) Update(t Executor, administratorAuthentication *AdministratorAuthentication) (bool, error) {
	count, err := t.Update(administratorAuthentication)
	return count > 0, err
}

func (r *administratorAuthenticationRepository) Delete(t Executor, administratorAuthentication *AdministratorAuthentication) (bool, error) {
	count, err := t.Delete(administratorAuthentication)
	return count > 0, err
}

func (r *authenticationRepository) FindByPersonID(t Executor, personID uint) ([]*Authentication, error) {
	var authentications []*Authentication

	_, err := t.Select(
		&authentications,
		"SELECT * FROM authentications WHERE person_id = :person_id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{
			"person_id":  personID,
			"account_id": r.AccountID,
		},
	)

	return authentications, err
}

func (r *authenticationRepository) FindByPersonIDAndProviderID(t Executor, personID uint, providerID uint) (*Authentication, error) {
	var (
		authentications []*Authentication
		authentication  *Authentication
	)

	_, err := t.Select(
		&authentications,
		"SELECT * FROM authentications WHERE person_id = :person_id AND provider_id = :provider_id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{
			"person_id":   personID,
			"provider_id": providerID,
			"account_id":  r.AccountID,
		},
	)

	if len(authentications) == 1 {
		authentication = authentications[0]
	}

	return authentication, err
}

func (r *authenticationRepository) FindByProviderIDAndToken(t Executor, providerID uint, token string) (*Authentication, error) {
	var (
		authentications []*Authentication
		authentication  *Authentication
	)

	_, err := t.Select(
		&authentications,
		"SELECT * FROM authentications WHERE token = :token AND provider_id = :provider_id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{
			"token":       token,
			"provider_id": providerID,
			"account_id":  r.AccountID,
		},
	)

	if len(authentications) == 1 {
		authentication = authentications[0]
	}

	return authentication, err
}

// Find an Authentication entity by it's provider, person and token values.
func (r *authenticationRepository) FindByProviderIDAndPersonIDAndToken(t Executor, providerID uint, personID uint, token string) (*Authentication, error) {
	var (
		authentications []*Authentication
		authentication  *Authentication
	)

	_, err := t.Select(
		&authentications,
		"SELECT * FROM authentications WHERE token = :token AND provider_id = :provider_id AND person_id = :person_id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{
			"token":       token,
			"person_id":   personID,
			"provider_id": providerID,
			"account_id":  r.AccountID,
		},
	)

	if len(authentications) == 1 {
		authentication = authentications[0]
	}

	return authentication, err
}

func (r *authenticationRepository) Create(t Executor, authentication *Authentication) error {
	authentication.AccountID = r.AccountID
	return t.Insert(authentication)
}

func (r *authenticationRepository) Update(t Executor, authentication *Authentication) (bool, error) {
	count, err := t.Update(authentication)
	return count > 0, err
}

func (r *authenticationRepository) Delete(t Executor, authentication *Authentication) (bool, error) {
	count, err := t.Delete(authentication)
	return count > 0, err
}

func (r *authenticationRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func (r *bridgeUserRepository) FindByPersonIDAndBridgeID(t Executor, pID uint, bID uint) (*BridgeUser, error) {
	var (
		user  *BridgeUser
		users []*BridgeUser
	)

	_, err := t.Select(
		&users,
		"SELECT * FROM bridge_users WHERE account_id = :account_id AND person_id = :person_id AND bridge_id = :bridge_id",
		map[string]interface{}{
			"account_id": r.AccountID,
			"person_id":  pID,
			"bridge_id":  bID,
		},
	)

	if len(users) == 1 {
		user = users[0]
	}

	return user, err
}

func (r *bridgeUserRepository) FindByPersonID(t Executor, pID uint) ([]*BridgeUser, error) {
	var users []*BridgeUser

	_, err := t.Select(
		&users,
		"SELECT * FROM bridge_users WHERE account_id = :account_id AND person_id = :person_id",
		map[string]interface{}{
			"account_id": r.AccountID,
			"person_id":  pID,
		},
	)

	return users, err
}

func (r *bridgeUserRepository) FindByBridgeID(t Executor, bid uint) ([]*BridgeUser, error) {
	var users []*BridgeUser

	_, err := t.Select(
		&users,
		"SELECT * FROM bridge_users WHERE account_id = :account_id AND bridge_id = :bridge_id",
		map[string]interface{}{
			"account_id": r.AccountID,
			"bridge_id":  bid,
		},
	)

	return users, err
}

func (r *bridgeUserRepository) Create(t Executor, bridgeUser *BridgeUser) error {
	bridgeUser.AccountID = r.AccountID
	return t.Insert(bridgeUser)
}

func (r *bridgeUserRepository) Update(t Executor, bridgeUser *BridgeUser) (bool, error) {
	count, err := t.Update(bridgeUser)
	return count > 0, err
}

func (r *bridgeUserRepository) Delete(t Executor, bridgeUser *BridgeUser) (bool, error) {
	count, err := t.Delete(bridgeUser)
	return count > 0, err
}

func (r *bridgeUserRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func (r *deviceRepository) All(t Executor) ([]*Device, error) {
	var devices []*Device

	_, err := t.Select(
		&devices,
		"SELECT * FROM devices WHERE account_id = :account_id ORDER BY NAME ASC",
		map[string]interface{}{"account_id": r.AccountID},
	)

	return devices, err
}

func (r *deviceRepository) Find(t Executor, id uint) (*Device, error) {
	var (
		devices []*Device
		device  *Device
	)

	_, err := t.Select(
		&devices,
		"SELECT * FROM devices WHERE id = :id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"id": id, "account_id": r.AccountID},
	)

	if len(devices) == 1 {
		device = devices[0]
	}

	return device, err
}

func (r *deviceRepository) FindByDeviceID(t Executor, deviceID string) (*Device, error) {
	var (
		devices []*Device
		device  *Device
	)

	_, err := t.Select(
		&devices,
		"SELECT * FROM devices WHERE device_id = :id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"device_id": deviceID, "account_id": r.AccountID},
	)

	if len(devices) == 1 {
		device = devices[0]
	}

	return device, err
}

func (r *deviceRepository) FindByToken(t Executor, token string) (*Device, error) {
	var (
		devices []*Device
		device  *Device
	)

	_, err := t.Select(
		&devices,
		"SELECT * FROM devices WHERE token = :token AND account_id := :account_id LIMIT 1",
		map[string]interface{}{"token": token, "account_id": r.AccountID},
	)

	if len(devices) == 1 {
		device = devices[0]
	}

	return device, err
}

func (r *deviceRepository) Create(t Executor, device *Device) error {
	device.AccountID = r.AccountID
	return t.Insert(device)
}

func (r *deviceRepository) Update(t Executor, device *Device) (bool, error) {
	count, err := t.Update(device)
	return count > 0, err
}

func (r *deviceRepository) Delete(t Executor, device *Device) (bool, error) {
	count, err := t.Delete(device)
	return count > 0, err
}

func (r *deviceRepository) Enable(t Executor, device *Device, enabled bool) (bool, error) {
	result, err := t.Exec(
		"UPDATE devices SET is_enabled = :enabled WHERE id = :id AND account_id = :account_id",
		map[string]interface{}{
			"enabled":    enabled,
			"id":         device.ID,
			"account_id": r.AccountID,
		},
	)

	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()

	return affected == 1, err
}

func (r *deviceRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func (r *doorRepository) All(t Executor) ([]*Door, error) {
	var doors []*Door

	_, err := t.Select(
		&doors,
		"SELECT * FROM doors WHERE account_id = :account_id ORDER BY name ASC",
		map[string]interface{}{"account_id": r.AccountID},
	)

	return doors, err
}

func (r *doorRepository) Find(t Executor, id uint) (*Door, error) {
	var (
		doors []*Door
		door  *Door
	)

	_, err := t.Select(
		&doors,
		"SELECT * FROM doors WHERE id = :id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"id": id, "account_id": r.AccountID},
	)

	if len(doors) == 1 {
		door = doors[0]
	}

	return door, err
}

func (r *doorRepository) Create(t Executor, door *Door) error {
	door.AccountID = r.AccountID
	return t.Insert(door)
}

func (r *doorRepository) Update(t Executor, door *Door) (bool, error) {
	count, err := t.Update(door)
	return count > 0, err
}

func (r *doorRepository) Delete(t Executor, door *Door) (bool, error) {
	count, err := t.Delete(door)
	return count > 0, err
}

func (r *doorRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func (r *eventRepository) All(t Executor, p *Pagination) ([]*Event, error) {
	var events []*Event

	parameters := map[string]interface{}{}

	query := "SELECT * FROM events"

	if r.AccountID > 0 {
		query += " WHERE account_id = :accountRepositoryunt_id"
		parameters["account_id"] = r.AccountID
	}

	query += fmt.Sprintf(" ORDER BY id %s LIMIT %d, %d", p.Order, p.Limit, (p.Page-1)*p.Limit)

	_, err := t.Select(
		&events,
		query,
		parameters,
	)

	return events, err
}

func (r *eventRepository) Find(t Executor, id uint) (*Event, error) {
	var (
		events []*Event
		event  *Event
	)
	parameters := map[string]interface{}{}

	query := "SELECT * FROM events WHERE id = :id LIMIT 1"

	if r.AccountID > 0 {
		query += " AND account_id = :account_id"
		parameters["account_id"] = r.AccountID
	}

	query += " LIMIT 1"

	_, err := t.Select(
		&events,
		query,
		parameters,
	)

	if len(events) == 1 {
		event = events[0]
	}

	return event, err
}

func (r *eventRepository) Create(t Executor, event *Event) error {
	event.AccountID = r.AccountID
	return t.Insert(event)
}

func (r *eventRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func (r *personRepository) All(t Executor) ([]*Person, error) {
	var people []*Person

	_, err := t.Select(
		&people,
		"SELECT * FROM people WHERE account_id = :account_id ORDER BY name ASC",
		map[string]interface{}{"account_id": r.AccountID},
	)

	return people, err
}

func (r *personRepository) Find(t Executor, id uint) (*Person, error) {
	var (
		people []*Person
		person *Person
	)

	_, err := t.Select(
		&people,
		"SELECT * FROM people WHERE id = :id AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"id": id, "account_id": r.AccountID},
	)

	if len(people) == 1 {
		person = people[0]
	}

	return person, err
}

func (r *personRepository) FindByEmail(t Executor, email string) (*Person, error) {
	var (
		people []*Person
		person *Person
	)

	_, err := t.Select(
		&people,
		"SELECT * FROM people WHERE email = :email AND account_id = :account_id LIMIT 1",
		map[string]interface{}{"email": email, "account_id": r.AccountID},
	)

	if len(people) == 1 {
		person = people[0]
	}

	return person, err
}

// Create a new person, setting the repository AccountID.
func (r *personRepository) Create(t Executor, person *Person) error {
	person.AccountID = r.AccountID
	return t.Insert(person)
}

func (r *personRepository) Update(t Executor, person *Person) (bool, error) {
	count, err := t.Update(person)
	return count > 0, err
}

func (r *personRepository) Delete(t Executor, person *Person) (bool, error) {
	count, err := t.Delete(person)
	return count > 0, err
}

func (r *personRepository) SetAccountScope(accountID uint) {
	r.AccountID = accountID
}

func initDatabase(c *DoorbotConfig) *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", c.Database.URL)

	if err != nil {
		panic(err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Account{}, "accounts").SetKeys(true, "ID")
	dbmap.AddTableWithName(Authentication{}, "authentications").SetKeys(true, "ID")
	dbmap.AddTableWithName(BridgeUser{}, "bridge_users").SetKeys(false)
	dbmap.AddTableWithName(Door{}, "doors").SetKeys(true, "ID")
	dbmap.AddTableWithName(Device{}, "devices").SetKeys(true, "ID")
	dbmap.AddTableWithName(Event{}, "events").SetKeys(true, "ID")
	dbmap.AddTableWithName(Person{}, "people").SetKeys(true, "ID")

	if c.Database.Trace {
		logger := log.New()
		dbmap.TraceOn("[gorp]", logger)
	}

	return dbmap
}

// MapDatabase maps the database in Martini
func MapDatabase(m *martini.Martini, config *DoorbotConfig) {
	m.Map(initDatabase(config))
}

// UseRepositories maps the repositories within a martini context.
func UseRepositories(m *martini.Martini) {

	m.Use(func(c martini.Context, database *gorp.DbMap) {
		c.MapTo(newRepositories(database), (*Repositories)(nil))
	})
}
