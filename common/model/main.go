package model

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

const (
	// User clearance level
	UserClearanceNone          = 0
	UserClearanceTeller        = 1
	UserClearanceAdministrator = math.MaxInt8

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

	// Service permissions
	ServicePermissionDebit     = 1 << 0
	ServicePermissionCredit    = 1 << 1
	ServicePermissionOverdraft = 1 << 2

	// Transaction state
	TransactionStateProcessing = "PRC"
	TransactionStateError      = "ERR"
	TransactionStateSuccess    = "SUC"
)

type User struct {
	Id        uuid.UUID
	Clearance int8
	Username  string
	Passhash  string
	Fullname  string
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

type Service struct {
	Id          uuid.UUID
	Type        string
	State       string
	Permissions int64
	Currency    string
	InitBalance decimal.Decimal
	Balance     decimal.Decimal
}

func NewService(mType, currency string, initBalance decimal.Decimal) (Service, error) {
	newService := Service{
		Type:        mType,
		State:       ServiceStateRequested,
		Permissions: ServicePermissionDebit | ServicePermissionCredit,
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

func (srv *Service) CheckPermissions(mask int64) bool {
	p := srv.Permissions & mask
	if p == 0 {
		return false
	} else {
		return true
	}
}

func (srv *Service) Debit(amount decimal.Decimal) error {
	if !srv.CheckPermissions(ServicePermissionDebit) {
		return fmt.Errorf("service %s does not have debit permission", srv.Id)
	}
	newBalance := srv.Balance.Sub(amount)
	if newBalance.Add(srv.InitBalance).IsNegative() &&
		!srv.CheckPermissions(ServicePermissionOverdraft) {
		return fmt.Errorf("service %s does not have overdraft permission", srv.Id)
	}
	srv.Balance = newBalance
	return nil
}

func (srv *Service) Credit(amount decimal.Decimal) error {
	if !srv.CheckPermissions(ServicePermissionDebit) {
		return fmt.Errorf("service %s does not have credit permission", srv.Id)
	}
	srv.Balance = srv.Balance.Add(amount)
	return nil
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
