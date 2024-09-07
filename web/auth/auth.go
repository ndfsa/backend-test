package main

import (
	"encoding/json"
	"net/http"

	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
	"github.com/ndfsa/cardboard-bank/web/token"
)

type AuthHandlerFactory struct {
	repo repository.AuthRepository
	mdf  middleware.MiddlewareFactory
}

func NewAuthHandlerFactory(
	repo repository.AuthRepository,
	mdf middleware.MiddlewareFactory,
) AuthHandlerFactory {
	return AuthHandlerFactory{repo, mdf}
}

func (factory *AuthHandlerFactory) Authenticate() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000))
	f := func(w http.ResponseWriter, r *http.Request) error {
		var req dto.AuthRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		user, err := factory.repo.Authenticate(r.Context(), req.Username, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}

		accessToken, err := token.GenerateAccessToken(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}
		refreshToken, err := token.GenerateRefreshToken(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}

		res := dto.AuthResponseDTO{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}
