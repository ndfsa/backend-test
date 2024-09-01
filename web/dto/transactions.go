package dto

import (
	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/shopspring/decimal"
)

type CreateTransactionRequestDTO struct {
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func (data *CreateTransactionRequestDTO) Parse() (model.Transaction, error) {
	amount, err := decimal.NewFromString(data.Amount)
	if err != nil {
		return model.Transaction{}, err
	}

	src, err := uuid.Parse(data.Source)
	if err != nil {
		return model.Transaction{}, err
	}

	dst, err := uuid.Parse(data.Destination)
	if err != nil {
		return model.Transaction{}, err
	}

	transaction, err := model.NewTransaction(data.Currency, amount, src, dst)
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

type CreateTransactionResponseDTO struct {
	Id string `json:"id"`
}

type ReadTransactionResponseDTO struct {
	Id          string
	State       string
	Time        string
	Currency    string
	Amount      string
	Source      string
	Destination string
}
