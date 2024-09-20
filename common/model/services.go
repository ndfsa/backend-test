package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
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
)

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
