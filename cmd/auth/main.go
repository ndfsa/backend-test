package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/cardboard-bank/cmd/auth/dto"
	"github.com/ndfsa/cardboard-bank/cmd/auth/repository"
	"github.com/ndfsa/cardboard-bank/internal/encoding"
	"github.com/ndfsa/cardboard-bank/internal/middleware"
	"github.com/ndfsa/cardboard-bank/internal/token"
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
	http.Handle("GET /auth", middleware.Basic(refreshToken(repo)))
	http.Handle("POST /auth/signup", middleware.Basic(signUp(repo)))

	// start server
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func authorizeUser(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user dto.AuthUserDTO
		if err := encoding.Receive(r, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if user.Username == "" || user.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("missing credentials")
			return
		}

		userId, err := repo.AuthenticateUser(r.Context(), user.Username, user.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		accessTokenString, err := token.GenerateAccessToken(userId, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		refreshTokenString, err := token.GenerateRefreshToken(userId, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		encoding.Send(w, dto.TokenDTO{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		})
	}
}

func signUp(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser dto.SignUpDTO
		if err := encoding.Receive(r, &newUser); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Println(err)
			return
		}

		err := repo.CreateUser(r.Context(), newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}

func refreshToken(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")
		if err := token.ValidateAccessToken(encodedToken, tokenKey); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		accessTokenString, err := token.GenerateAccessToken(userId, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		refreshTokenString, err := token.GenerateRefreshToken(userId, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		encoding.Send(w, dto.TokenDTO{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		})
	}
}
