package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

const (
	// User role
	UserRoleRegular       = "USR"
	UserRoleOfficer       = "OFC"
	UserRoleAdministrator = "ADM"

	// Currency
	CurrencyUnitedStatesDollar = "USD"
	CurrencyCanadianDollar     = "CAD"
	CurrencyJapaneseYen        = "JPY"
	CurrencyNorwegianCrown     = "NOK"

	// Service type
	ServiceTypeSavings              = "SAV"
	ServiceTypeChequing             = "CHQ"
	ServiceTypeLoan                 = "LOA"
	ServiceTypeLineOfCredit         = "LOC"
	ServiceTypeCertificateOfDeposit = "COD"

	// Service state
	ServiceStateRequested = "REQ"
	ServiceStateActive    = "ACT"
	ServiceStateFrozen    = "FRZ"
	ServiceStateClosed    = "CLS"

	// Transaction state
	TransactionStateProcessing = "PRC"
	TransactionStateError      = "ERR"
	TransactionStateSuccess    = "SUC"
)

type User struct {
	Id       uuid.UUID
	Role     string
	Username string
	Passhash string
	Fullname string
}

func (user *User) Validate(password string) error {
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
		Role:     UserRoleRegular,
		Username: username,
		Fullname: fullname,
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

type Service struct {
	Id          uuid.UUID
	Type        string
	State       string
	Currency    string
	InitBalance decimal.Decimal
	Balance     decimal.Decimal
}

func NewService(mType, currency string, initBalance decimal.Decimal) (Service, error) {
	newService := Service{
		Type:        mType,
		State:       ServiceStateRequested,
		Currency:    currency,
		InitBalance: initBalance,
		Balance:     decimal.Zero,
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Service{}, err
	}
	newService.Id = id

	return newService, nil
}

type Transaction struct {
	Id          uuid.UUID
	State       string
	Time        string
	Currency    string
	Amount      decimal.Decimal
	Source      uuid.UUID
	Destination uuid.UUID
}

func NewTransaction(
	currency string,
	amount decimal.Decimal,
	src,
	dst uuid.UUID,
) (Transaction, error) {
	newTransaction := Transaction{
		State:       TransactionStateProcessing,
		Time:        "NOW",
		Currency:    currency,
		Amount:      amount,
		Source:      src,
		Destination: dst,
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Transaction{}, err
	}
	newTransaction.Id = id

	return newTransaction, nil
}
