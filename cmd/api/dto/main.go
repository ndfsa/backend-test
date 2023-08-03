package dto

type UserDto struct {
	Fullname string `json:"name"`
	Username string `json:"user"`
	Password string `json:"pass"`
}
