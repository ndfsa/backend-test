package token

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	User uint64 `json:"user"`
	jwt.RegisteredClaims
}

func GetUserId(bearerToken string) (uint64, error) {
	tokenString := strings.Split(bearerToken, " ")[1]
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &CustomClaims{})
    if err != nil {
        return 0, err
    }
	claims, ok := token.Claims.(*CustomClaims)
    if !ok {
        return 0, errors.New("failed casting claims")
    }
	return claims.User, nil
}

func Validate(header string) error {
	tokenString := strings.TrimSpace(header)
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
