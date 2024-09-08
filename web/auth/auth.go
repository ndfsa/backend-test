package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
	"github.com/ndfsa/cardboard-bank/web/token"
)

type AuthHandlerFactory struct {
	repo repository.AuthRepository
}

func NewAuthHandlerFactory(
	repo repository.AuthRepository,
) AuthHandlerFactory {
	return AuthHandlerFactory{repo}
}

func (factory *AuthHandlerFactory) Authenticate() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000))
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.AuthRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		user, err := factory.repo.Authenticate(r.Context(), req.Username, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		accessToken, err := token.GenerateAccessToken(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		refreshToken, err := token.GenerateRefreshToken(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		res := dto.AuthResponseDTO{
			Id:           user.Id.String(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}
