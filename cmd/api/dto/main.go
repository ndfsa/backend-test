package dto

import "github.com/shopspring/decimal"

type UserDto struct {
	Fullname string `json:"name"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

type ServiceDto struct {
	Type        string          `json:"type"`
	Currency    string          `json:"currency"`
	InitBalance decimal.Decimal `json:"init_balance"`
}
