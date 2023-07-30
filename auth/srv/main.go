package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ndfsa/backend-test/middleware"
)

var users map[string]uint64

func main() {
	users = make(map[string]uint64)

	users["root"] = 0

	http.Handle("/auth", middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000))(http.HandlerFunc(auth)))

	err := http.ListenAndServe(":4001", nil)
	log.Fatal(err)
}

func generateJWT(userId uint64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userId,
		"exp":  time.Now().Add(10 * time.Second).Unix(),
		"nbf":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte("test-application"))
	if err != nil {
		log.Printf("error generating jwt: %s\n", err.Error())
	}

	return tokenString
}

type UserDTO struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func auth(w http.ResponseWriter, r *http.Request) {
	var user UserDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("authentication: error %s\n", err.Error())
		http.Error(w, "malformed json", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		log.Println("authentication: error missing parameters")
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	// pull user from storage
	userId, ok := users[user.Username]
	if !ok ||
        // mock password validation
		user.Username != user.Password {

		log.Printf("authentication: no such username/password=%s/****\n",
			user.Username)
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "No such user", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["token"] = generateJWT(userId)

	json.NewEncoder(w).Encode(resp)
	log.Printf("authentication: user authenitcation successful for user=%s\n", user.Username)
}
