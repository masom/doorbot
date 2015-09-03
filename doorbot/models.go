package doorbot

import (
	"time"
)

// Account represent top-level entities under which doors and people are grouped.
type Account struct {
	ID        uint   `db:"id" json:"id"`
	Name      string `db:"name" form:"name" json:"name" binding:"required"`
	IsEnabled bool   `db:"is_enabled" form:"is_enabled" json:"is_enabled"`
	Host      string `db:"host" form:"host" json:"host"`

	BridgeHubEnabled bool   `db:"bridge_hub_enabled" json:"bridge_hub_enabled"`
	BridgeHubURL     string `db:"bridge_hub_url" json:"bridge_hub_url"`
	BridgeHubToken   string `db:"bridge_hub_token" json:"bridge_hub_token"`

	BridgeHipChatEnabled bool `db:"bridge_hipchat_enabled" json:"bridge_hipchat_enabled"`
	BridgeHipChatToken string `db:"bridge_hipchat_token" json:"bridge_hipchat_token"`

	BridgeSlackEnabled bool `db:"bridge_slack_enabled" json:"bridge_slack_enabled"`
	BridgeSlackToken string `db:"bridge_slack_token" json:"bridge_slack_token"`

	ContactName           string `db:"contact_name" json:"contact_name"`
	ContactEmail          string `db:"contact_email" json:"contact_email"`
	ContactEmailConfirmed bool   `db:"contact_email_confirmed" json:"-"`
	ContactPhoneNumber    string `db:"contact_phone_number" json:"contact_phone_number"`

	NotificationsEnabled bool `db:"notifications_enabled" json:"notifications_enabled"`

	NotificationsEmailMessageTemplate *string `db:"notifications_email_message_template" json:"notifications_email_message_template"`

	NotificationsHipChatEnabled bool `db:"notifications_hipchat_enabled" json:"notifications_hipchat_enabled"`
	NotificationsHipChatToken string `db:"notifications_hipchat_token" json:"notifications_hipchat_token"`

	NotificationsMailgunEnabled         bool    `db:"notifications_mailgun_enabled" json:"notifications_mailgun_enabled"`

	NotificationsNexmoEnabled bool `db:"notifications_nexmo_enabled" json:"notifications_nexmo_enabled"`
	NotificationsNexmoToken string `db:"notifications_nexmo_token" json:"notifications_nexmo_token"`

	NotificationsPostmarkEnabled         bool    `db:"notifications_postmark_enabled" json:"notifications_postmark_enabled"`

	NotificationsSlackEnabled bool `db:"notifications_slack_enabled" json:"notifications_slack_enabled"`
	NotificationsSlackToken string `db:"notifications_slack_token" json:"notifications_slack_token"`

	NotificationsSMSMessageTemplate *string `db:"notifications_sms_message_template" json:"notifications_sms_message_template"`

	NotificationsTwilioEnabled           bool    `db:"notifications_twilio_enabled" json:"notifications_twilio_enabled"`
	NotificationsTwilioSourcePhoneNumber *string `db:"notifications_twilio_source_phone_number" json:"notifications_twilio_source_phone_number"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Administrator holds data about Doorbot super users
type Administrator struct {
	ID    uint   `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
	Name  string `db:"name" json:"name"`
}

// AdministratorAuthentication holds authentication data for adminnistrators
type AdministratorAuthentication struct {
	AdministratorID uint      `db:"administrator_id" json:"administrator_id"`
	ProviderID      uint      `db:"provider_id" json:"provider_id"`
	Token           string    `db:"token" json:"_"`
	LastUsedAt      string    `db:"last_used_at" json:"last_used_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Authentication represent authentication data ( oauth, password, etc.) attached to a specific Account
type Authentication struct {
	ID         uint   `db:"id"`
	AccountID  uint   `db:"account_id"`
	PersonID   uint   `db:"person_id"`
	ProviderID uint   `db:"provider_id"`
	Token      string `db:"token"`
}

// BridgeUser represents the link between a bridge user and a doorbot person
type BridgeUser struct {
	AccountID uint      `db:"account_id"`
	PersonID  uint      `db:"person_id"`
	BridgeID  uint      `db:"bridge_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Name        string `db:"-"`
	Email       string `db:"-"`
	PhoneNumber string `db:"-"`
	Title       string `db:"-"`
}

// Door ... The door app client will report itself as one of these.
type Door struct {
	AccountID uint   `db:"account_id" json:"account_id"`
	ID        uint   `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
}

// Event holds event data.
type Event struct {
	ID        uint      `db:"event_id" json:"event_id"`
	AccountID uint      `db:"account_id" json:"account_id"`
	DoorID    uint      `db:"door_id" json:"door_id"`
	DeviceID  uint      `db:"device_id" json:"device_id"`
	EventID   uint      `db:"event_id" json:"event_id"`
	PersonID  uint      `db:"person_id" json:"person_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Device holds device data
type Device struct {
	AccountID   uint      `db:"account_id" json:"-"`
	ID          uint      `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	DeviceID    *string   `db:"device_id" json:"device_id"`
	DoorID      *uint     `db:"door_id" json:"door_id"`
	Make        string    `db:"make" json:"make"`
	Description string    `db:"description" json:"description"`
	IsEnabled   bool      `db:"is_enabled" json:"is_enabled"`
	Token       string    `db:"token" json:"token"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Person is someone who can be reached behind a Door
type Person struct {
	AccountID   uint   `db:"account_id" json:"account_id"`
	AccountType uint   `db:"account_type" json:"account_type"`
	ID          uint   `db:"id" json:"id"`
	Name        string `db:"name" form:"name" json:"name" binding:"required"`
	Title       string `db:"title" form:"title" json:"title"`
	Email       string `db:"email" form:"email" json:"email" binding:"required"`
	PhoneNumber string `db:"phone_number" form:"phone_number" json:"phone_number"`
	IsVisible   bool   `db:"is_visible" form:"is_visible" json:"is_visible"`
	IsAvailable bool   `db:"is_available" form:"is_available" json:"is_available"`

	NotificationsEnabled      bool `db:"notifications_enabled" form:"notifications_enabled" json:"notifications_enabled"`
	NotificationsAppEnabled   bool `db:"notifications_app_enabled" form:"notifications_app_enabled" json:"notifications_app_enabled"`
	NotificationsEmailEnabled bool `db:"notifications_email_enabled" form:"notifications_email_enabled" json:"notifications_email_enabled"`
	NotificationsChatEnabled   bool `db:"notifications_chat_enabled" form:"notifications_chat_enabled" json:"notifications_chat_enabled"`
	NotificationsSMSEnabled   bool `db:"notifications_sms_enabled" form:"notifications_sms_enabled" json:"notifications_sms_enabled"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// IsAccountManager determine if the person is an account manager
func (p *Person) IsAccountManager() bool {
	return p.AccountType == AccountOwner || p.AccountType == AccountManager
}

// IsAccountOwner determine if the person is the account owner
func (p *Person) IsAccountOwner() bool {
	return p.AccountType == AccountOwner
}

// PersonArguments represents the different arguments required when creating a new Person
type PersonArguments struct {
	AccountID   uint
	Name        string
	Title       string
	Blurb       string
	Email       string
	PhoneNumber string
	IsVisible   bool
}

// NewAuthentication creates a new Authentication
func NewAuthentication(accountID uint, providerID uint, token string) *Authentication {
	return &Authentication{
		AccountID:  accountID,
		ProviderID: providerID,
		Token:      token,
	}
}

// NewDoor creates a new Door
func NewDoor(accountID uint, name string) *Door {
	return &Door{
		AccountID: accountID,
		Name:      name,
	}
}

// NewPerson creates a new Person
func NewPerson(args PersonArguments) *Person {
	return &Person{
		AccountID:   args.AccountID,
		Name:        args.Name,
		Title:       args.Title,
		Email:       args.Email,
		PhoneNumber: args.PhoneNumber,
		IsVisible:   args.IsVisible,
	}
}

// NewEvent creates a new event
func NewEvent(accountID uint, deviceID uint, doorID uint, eventID uint, personID uint) *Event {
	return &Event{
		AccountID: accountID,
		DeviceID:  deviceID,
		DoorID:    doorID,
		EventID:   eventID,
		PersonID:  personID,
	}
}
