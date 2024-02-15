package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/backend-test/cmd/auth/dto"
	"github.com/ndfsa/backend-test/cmd/auth/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

const (
    // 32-byte key
	tokenKey = "2bbb515c1311dd69a609a0d553dc7ac1ac8eadc2b22daa9aaa99483d2f381374"
)

func main() {
	// connect to database
	db, err := sql.Open("pgx", "postgres://back:root@db:5432/cardboard_bank")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

    log.Println(tokenKey)

	repo := repository.NewAuthRepository(db)

	// setup routes
	http.Handle("POST /auth", middleware.Basic(authorizeUser(repo)))
	http.Handle("POST /auth/signup", middleware.Basic(signUp(repo)))
	http.Handle("POST /auth/refresh", middleware.Basic(refreshToken(repo)))

	// start server
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func authorizeUser(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user dto.AuthUserDTO
		if err := util.Receive(r.Body, &user); err != nil {
			util.SendError(&w, http.StatusBadRequest, err.Error())
			return
		}

		if user.Username == "" || user.Password == "" {
			util.SendError(&w, http.StatusBadRequest, "missing credentials")
			return
		}

		userId, err := repo.AuthenticateUser(r.Context(), user.Username, user.Password)
		if err != nil {
			util.SendError(&w, http.StatusUnauthorized, err.Error())
			return
		}

		accessTokenString, err := token.GenerateAccessToken(userId, tokenKey)
		if err != nil {
			util.SendError(&w, http.StatusUnauthorized, err.Error())
			return
		}
		refreshTokenString, err := token.GenerateRefreshToken(userId, tokenKey)
		if err != nil {
			util.SendError(&w, http.StatusUnauthorized, err.Error())
			return
		}

		util.Send(&w, dto.TokenDTO{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		})
	}
}

func signUp(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.CreateUser(r.Context(), r.Body)
		if err != nil {
			util.SendError(&w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func refreshToken(repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate tokens

		// generate new tokens

		// send new tokens
	}
}
