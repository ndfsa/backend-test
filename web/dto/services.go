package dto

import (
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

type ReadServiceResponseDTO struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	State       string `json:"state"`
	Currency    string `json:"currency"`
	InitBalance string `json:"init_balance"`
	Balance     string `json:"balance"`
}
