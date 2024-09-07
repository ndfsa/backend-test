package dto

import (
	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/shopspring/decimal"
)

type CreateServiceRequestDTO struct {
	Owner       string `json:"user_id"`
	Type        string `json:"type"`
	Currency    string `json:"currency"`
	InitBalance string `json:"init_balance"`
}

func (dto *CreateServiceRequestDTO) Parse() (model.Service, uuid.UUID, error) {
	initBalance, err := decimal.NewFromString(dto.InitBalance)
	if err != nil {
		return model.Service{}, uuid.UUID{}, err
	}

	owner, err := uuid.Parse(dto.Owner)
	if err != nil {
		return model.Service{}, uuid.UUID{}, err
	}

	service, err := model.NewService(dto.Type, dto.Currency, initBalance)
	if err != nil {
		return model.Service{}, uuid.UUID{}, err
	}

	return service, owner, nil
}

type CreateServiceResponseDTO struct {
	Id string `json:"id"`
}

type ReadServiceResponseDTO struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	State       string `json:"state"`
	Currency    string `json:"currency"`
	InitBalance string `json:"init_balance"`
	Balance     string `json:"balance"`
}

type UpdateServiceRequestDTO struct {
	Id      string `json:"id"`
	State   string `json:"state"`
	Balance string `json:"balance"`
}
