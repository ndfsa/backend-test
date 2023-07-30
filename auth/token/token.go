package token

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Validate(w http.ResponseWriter, r *http.Request) error {

	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

	if tokenString == "" {
		return errors.New("bearer token not found")
	}

	bearerToken := strings.Split(tokenString, " ")

	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return errors.New("invalid bearer token")
	}

	tokenString = bearerToken[1]

	token, err := jwt.Parse(tokenString,
		tokenValidator,
		jwt.WithLeeway(1*time.Second),
		jwt.WithValidMethods([]string{"RS256", "HS256"}))

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid jwt token")
	}

	return nil
}

func tokenValidator(token *jwt.Token) (interface{}, error) {
	return []byte("test-application"), nil
}
