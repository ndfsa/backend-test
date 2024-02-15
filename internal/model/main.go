package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	UserId   uuid.UUID `json:"id"`
	Fullname string    `json:"fullname"`
	Username string    `json:"username"`
}

type Service struct {
	Id          uuid.UUID       `json:"id"`
	Type        string          `json:"type"`
	State       string          `json:"state"`
	Currency    string          `json:"currency"`
	InitBalance decimal.Decimal `json:"init_balance"`
	Balance     decimal.Decimal `json:"balance"`
}
