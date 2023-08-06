package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/backend-test/cmd/auth/dto"
	"github.com/ndfsa/backend-test/cmd/auth/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

func main() {
	// connect to database
	db, err := sql.Open("pgx", "postgres://back:root@localhost:5432/back_test")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// setup routes
	http.Handle("/auth/signup", middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000))(signUpHandler(db)))
	http.Handle("/auth", middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000))(auth(db)))

	// start server
	if err := http.ListenAndServe(":4001", nil); err != nil {
		log.Fatal(err)
	}
}

func generateJWT(userId uint64) (string, error) {
	claims := token.CustomClaims{
		User: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("test-application"))
	if err != nil {
        return "", err
	}

	return tokenString, nil
}

func auth(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user dto.AuthUserDTO
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}

		if user.Username == "" || user.Password == "" {
            util.Error(&w, http.StatusBadRequest, "missing credentials")
			return
		}

		userId, err := repository.AuthenticateUser(db, user.Username, user.Password)
		if err != nil {
            util.Error(&w, http.StatusUnauthorized, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]string)
		resp["token"], err = generateJWT(userId)
        if err != nil {
            util.Error(&w, http.StatusUnauthorized, err.Error())
            return
        }

		json.NewEncoder(w).Encode(resp)
	})
}

func signUpHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := repository.SignUp(db, r.Body)
		if err != nil {
            util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Fprint(w, userId)
	})
}
