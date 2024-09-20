package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	// Transaction state
	TransactionStateProcessing = "PRC"
	TransactionStateError      = "ERR"
	TransactionStateSuccess    = "SUC"
)

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
