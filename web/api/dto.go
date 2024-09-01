package main

import (
	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/shopspring/decimal"
)

type CreateServiceRequestDTO struct {
	Type        string `json:"type"`
	Currency    string `json:"currency"`
	InitBalance string `json:"init_balance"`
}

func (dto *CreateServiceRequestDTO) Parse() (model.Service, error) {
	initBalance, err := decimal.NewFromString(dto.InitBalance)
	if err != nil {
		return model.Service{}, err
	}

	service, err := model.NewService(dto.Type, dto.Currency, initBalance)
	if err != nil {
		return model.Service{}, err
	}

	return service, nil
}

type CreateServiceResponseDTO struct {
	Id string `json:"id"`
}

type CreateUserRequestDTO struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (dto *CreateUserRequestDTO) Parse() (model.User, error) {
	user, err := model.NewUser(dto.Username, dto.Fullname, dto.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

type CreateUserResponseDTO struct {
	Id string `json:"id"`
}

type ReadUserResponseDTO struct {
	Id       string `json:"id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}

type UpdateUserRequestDTO struct {
	Id       string `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (data *UpdateUserRequestDTO) Parse() (model.User, error) {
	userId, err := uuid.Parse(data.Id)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Id:       userId,
		Username: data.Username,
		Fullname: data.Fullname,
	}

	if data.Password == "" {
		return user, nil
	}

	if err := user.SetPassword(data.Password); err != nil {
		return model.User{}, err
	}

	return user, nil
}

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
