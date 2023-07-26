package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var users map[string]string

func main() {
	users = make(map[string]string)

	rootUserId := md5.Sum([]byte("root|toor"))
	users["root"] = hex.EncodeToString(rootUserId[:])

	http.Handle("/hello", authMiddleware(http.HandlerFunc(hello)))
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		bearerToken := strings.Split(tokenString, " ")
		if len(bearerToken) < 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString = bearerToken[1]

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("test-application"), nil
			}, jwt.WithLeeway(1*time.Second))

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized not valid", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST is supported", http.StatusBadRequest)
		log.Printf("authentication: received request with %s method\n", r.Method)
		return
	}
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
