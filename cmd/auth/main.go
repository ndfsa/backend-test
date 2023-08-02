package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
)

func main() {
	http.Handle("/auth", middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000))(http.HandlerFunc(auth)))

	err := http.ListenAndServe(":4001", nil)
	log.Fatal(err)
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

    // TODO authenticate user

	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["token"] = generateJWT(0)

	json.NewEncoder(w).Encode(resp)
	log.Printf("authentication: user authenitcation successful for user=%s\n", user.Username)
}
