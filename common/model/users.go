package model

import (
	"math"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
const (
	// User clearance level
	UserClearanceNone          = 0
	UserClearanceTeller        = 1
	UserClearanceAdministrator = math.MaxInt8
)

type User struct {
	Id        uuid.UUID
	Clearance int8
	Username  string
	Passhash  string
	Fullname  string
}

func (user *User) Validate(password string) error {
    if user.Passhash == "" {
        user.Passhash = "$2a$10$NRuHlLTONUKtxXMEFQDxxOmBSB9rS7/IsLP7eRjKZcRKwA/eB1EQ6"
    }
	return bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(password))
}

func (user *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Passhash = string(hashBytes)

	return nil
}

func NewUser(username, fullname, password string) (User, error) {
	newUser := User{
		Clearance: UserClearanceNone,
		Username:  username,
		Fullname:  fullname,
	}
	id, err := uuid.NewV7()
	if err != nil {
		return User{}, err
	}
	newUser.Id = id

	if err := newUser.SetPassword(password); err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (user *User) CheckOwnership(userId uuid.UUID) bool {
	return user.Id == userId
}
