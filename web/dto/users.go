package dto

import (
	"github.com/ndfsa/cardboard-bank/common/model"
)

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
	Id        string `json:"id"`
	Clearance int8   `json:"clearance"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
}

type UpdateUserRequestDTO struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (data *UpdateUserRequestDTO) Parse() (model.User, error) {
	user := model.User{
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
