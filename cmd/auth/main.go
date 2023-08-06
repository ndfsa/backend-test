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

func generateJWT(userId uint64) string {
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
		log.Printf("error generating jwt: %s\n", err.Error())
	}

	return tokenString
}

func auth(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user dto.AuthUserDTO
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "malformed json", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" {
			log.Println("user or password missing")
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		}

		userId, err := repository.AuthenticateUser(db, user.Username, user.Password)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "could not authenticate", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]string)
		resp["token"] = generateJWT(userId)

		json.NewEncoder(w).Encode(resp)
		log.Printf("authentication: user authenitcation successful for user=%s\n", user.Username)
	})
}

func signUpHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := repository.SignUp(db, r.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, userId)
	})
}
