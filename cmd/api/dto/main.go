package dto

import "github.com/shopspring/decimal"

type UpdateUserRequest struct {
	Fullname string `json:"name"`
	Username string `json:"user"`
}

type CreateServiceRequest struct {
	Type        string          `json:"type"`
	Currency    string          `json:"currency"`
	InitBalance decimal.Decimal `json:"init_balance"`
}

type CreateServiceResponse struct {
	Id string `json:"id"`
}

type ExecuteTransactionRequest struct {
	From     uint64          `json:"from"`
	To       uint64          `json:"to"`
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}
