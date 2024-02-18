package dto

type AuthUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenResponse struct {
    AccessToken string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
}
