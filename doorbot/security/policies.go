package security

// Policy represents CRUD access rights
type Policy struct {
	AccountCreate bool `json:"account_create,omitempty"`
	AccountDelete bool `json:"account_delete,omitempty"`
	AccountUpdate bool `json:"account_edit,omitempty"`
	AccountList   bool `json:"account_list,omitempty"`
	AccountView   bool `json:"account_view,omitempty"`

	DeviceCreate bool `json:"device_create,omitempty"`
	DeviceDelete bool `json:"device_delete,omitempty"`
	DeviceUpdate bool `json:"device_edit,omitempty"`
	DeviceList   bool `json:"device_list,omitempty"`
	DeviceView   bool `json:"device_view,omitempty"`

	DoorCreate bool `json:"door_create,omitempty"`
	DoorDelete bool `json:"door_delete,omitempty"`
	DoorUpdate bool `json:"door_edit,omitempty"`
	DoorList   bool `json:"door_list,omitempty"`
	DoorView   bool `json:"door_view,omitempty"`

	NotificationCreate bool `json:"notification_create,omitempty"`

	PersonCreate bool `json:"person_create,omitempty"`
	PersonDelete bool `json:"person_delete,omitempty"`
	PersonUpdate bool `json:"person_edit,omitempty"`
	PersonList   bool `json:"person_list,omitempty"`
	PersonView   bool `json:"person_view,omitempty"`
}

// EntityPolicy represents policy access to specific entities.
type EntityPolicy struct {
	CanUpdate bool `json:"can_edit,omitempty"`
	CanDelete bool `json:"can_delete,omitempty"`
	CanView   bool `json:"can_view,omitempty"`
}

// CanUpdatePerson determine if the current policy can edit a person OR if the two provided ids matches ( editing self )
func (p *Policy) CanUpdatePerson(currentPersonID uint, otherPersonID uint) bool {
	return p.PersonUpdate || currentPersonID == otherPersonID
}

// NewAdministratorPolicy creates an administrator policy
func NewAdministratorPolicy() *Policy {
	return &Policy{
		AccountCreate: true,
		AccountDelete: true,
		AccountUpdate: true,
		AccountList:   true,
		AccountView:   true,

		DeviceCreate: true,
		DeviceDelete: true,
		DeviceUpdate: true,
		DeviceList:   true,
		DeviceView:   true,

		DoorCreate: true,
		DoorDelete: true,
		DoorUpdate: true,
		DoorList:   true,
		DoorView:   true,

		PersonCreate: true,
		PersonDelete: true,
		PersonUpdate: true,
		PersonList:   true,
		PersonView:   true,
	}
}

// NewDevicePolicy creates a policy with device rights
func NewDevicePolicy() *Policy {
	return &Policy{
		DoorList: true,
		DoorView: true,

		PersonList: true,
	}
}

// NewOwnerPolicy creates a policy with owner rights
func NewOwnerPolicy() *Policy {
	return &Policy{
		AccountDelete: true,
		AccountUpdate: true,
		AccountView:   true,

		DeviceCreate: true,
		DeviceDelete: true,
		DeviceUpdate: true,
		DeviceList:   true,
		DeviceView:   true,

		DoorCreate: true,
		DoorDelete: true,
		DoorUpdate: true,
		DoorList:   true,
		DoorView:   true,

		PersonCreate: true,
		PersonDelete: true,
		PersonUpdate: true,
		PersonList:   true,
		PersonView:   true,
	}
}

// NewManagerPolicy creates a policy with manager rights
func NewManagerPolicy() *Policy {
	return &Policy{
		AccountUpdate: true,
		AccountView:   true,

		DeviceCreate: true,
		DeviceDelete: true,
		DeviceUpdate: true,
		DeviceList:   true,
		DeviceView:   true,

		DoorCreate: true,
		DoorDelete: true,
		DoorUpdate: true,
		DoorList:   true,
		DoorView:   true,

		PersonCreate: true,
		PersonDelete: true,
		PersonUpdate: true,
		PersonList:   true,
		PersonView:   true,
	}
}

// NewMemberPolicy creates a new policy with member rights
func NewMemberPolicy() *Policy {
	return &Policy{
		AccountView: true,

		DoorList: true,
		DoorView: true,

		PersonList: true,
		PersonView: true,
	}
}
