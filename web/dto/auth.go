package dto

type AuthRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponseDTO struct {
	Id           string `json:"id"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
