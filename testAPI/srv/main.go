package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ndfsa/backend-test/middleware"
)

var users map[string]string

func main() {
	users = make(map[string]string)

	rootUserId := md5.Sum([]byte("root|toor"))
	users["root"] = hex.EncodeToString(rootUserId[:])

	http.Handle("/user", middleware.Chain(
		middleware.Logger,
		middleware.Methods("POST", "GET"),
        middleware.UploadLimit(0),
		middleware.Auth)(http.HandlerFunc(getUser)))

	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]interface{})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		resp["user"] = claims["user"]
	} else {

	}

	json.NewEncoder(w).Encode(resp)
}
