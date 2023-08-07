package dto

// REQUEST
type AuthUserDTO struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type SignUpDTO struct {
	Fullname string `json:"name"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

// RESPONSE
type TokenDTO struct {
    Token string `json:"token"`
}
