package model

import (
	"errors"
	"net/mail"
	"time"
)

type User struct {
	ID              int64  `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Active          bool   `json:"active"`
	EmailVerified   bool   `json:"email_verified"`
	EmailSubscribed bool   `json:"email_subscribed"`
	Admin           bool   `json:"admin"`
	// Move to user monitoring service?
	Created    time.Time `json:"created"`
	LastLogin  time.Time `json:"last_login"`
	LastAction time.Time `json:"last_action"`
}

func NewUser(email string, password string, admin bool, now time.Time) (*User, error) {
	u := &User{
		Email:           email,
		Password:        password,
		Active:          true,
		EmailVerified:   false,
		EmailSubscribed: true,
		Admin:           admin,
		Created:         now,
		LastLogin:       now,
		LastAction:      now,
	}

	if err := u.validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) validate() error {
	if err := u.checkEmptyFields(); err != nil {
		return err
	}
	if err := u.checkEmail(); err != nil {
		return err
	}

	return nil
}

func (u *User) checkEmptyFields() error {
	if u.Password == "" {
		return errors.New("a required field is empty")
	}
	return nil
}

func (u *User) checkEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}
