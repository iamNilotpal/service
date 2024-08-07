package user

import (
	"net/mail"
	"time"
)

// User represents information about an individual user.
type User struct {
	ID           int
	FirstName    string
	LastName     string
	Email        mail.Address
	Enabled      bool
	Department   string
	PasswordHash []byte
	Roles        []Role
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string
	Email           mail.Address
	Roles           []Role
	Department      string
	Password        string
	PasswordConfirm string
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name            *string
	Email           *mail.Address
	Roles           []Role
	Department      *string
	Password        *string
	PasswordConfirm *string
	Enabled         *bool
}
