package main

type AuthRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
