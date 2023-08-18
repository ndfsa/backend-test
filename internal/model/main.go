package model

import "github.com/shopspring/decimal"

type User struct {
	UserId   uint64 `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
}

type Service struct {
	Id            uint64
	Type          uint8
	State         uint8
	InitBalance   decimal.Decimal
	DebitBalance  decimal.Decimal
	CreditBalance decimal.Decimal
}
