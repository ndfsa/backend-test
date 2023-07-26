package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var users map[string]string

func main() {
	users = make(map[string]string)

	rootUserId := md5.Sum([]byte("root|toor"))
	users["root"] = hex.EncodeToString(rootUserId[:])

	http.HandleFunc("/auth", auth)
	err := http.ListenAndServe(":4001", nil)
	log.Fatal(err)
}

func generateJWT(userId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userId,
		"exp": time.Now().Add(10 * time.Second).Unix(),
		"nbf": time.Now().Unix(),
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
	if r.Method != "POST" {
		http.Error(w, "Only POST is supported", http.StatusBadRequest)
		log.Printf("authentication: received request with %s method\n", r.Method)
		return
	}

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

	userId, ok := users[user.Username]
	if !ok {
		log.Printf("authentication: no such username/password=%s/****\n",
			user.Username)
		http.Error(w, "No such user", http.StatusForbidden)
		return
	}

	hash := md5.Sum([]byte(user.Username + "|" + user.Password))

	calculatedId := hex.EncodeToString(hash[:])
	if calculatedId != userId {
		log.Printf("authentication: user authenitcation failed for user=%s\n",
			user.Username)
		http.Error(w, "No such user", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["token"] = generateJWT(userId)

	json.NewEncoder(w).Encode(resp)
	log.Printf("authentication: user authenitcation successful for user=%s\n", user.Username)
}
