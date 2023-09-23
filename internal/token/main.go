package token

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	User uint64 `json:"user"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	User    uint64 `json:"user"`
	Refresh bool   `json:"refresh"`
	jwt.RegisteredClaims
}

func GetUserId(r *http.Request) (uint64, error) {
	bearerToken := r.Header.Get("Authorization")
	tokenString := strings.Split(bearerToken, " ")[1]
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &AccessTokenClaims{})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return 0, errors.New("failed casting claims")
	}
	return claims.User, nil
}

func GetUserIdRef(r *http.Request) (uint64, error) {
	bearerToken := r.Header.Get("Authorization")
	tokenString := strings.Split(bearerToken, " ")[1]
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &RefreshTokenClaims{})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return 0, errors.New("failed casting claims")
	}
	return claims.User, nil
}

func ValidateAccessToken(r *http.Request) error {
	tokenString, err := ValidateToken(r)
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(tokenString,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("test-application"), nil
		},
		jwt.WithLeeway(1*time.Second),
		jwt.WithValidMethods([]string{"RS256", "HS256"}))
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid jwt token")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return errors.New("invalid jwt claims")
	}
	if claims.User == 0 {
		return errors.New("invalid user")
	}

	return nil
}

func ValidateRefreshToken(r *http.Request) error {
	tokenString, err := ValidateToken(r)
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(tokenString,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("refresh-test-application"), nil
		},
		jwt.WithLeeway(5*time.Second),
		jwt.WithValidMethods([]string{"RS256", "HS256"}))
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid jwt token")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return errors.New("invalid jwt claims")
	}
	if claims.User == 0 {
		return errors.New("invalid user")
	}
	if !claims.Refresh {
		return errors.New("not a refresh token")
	}

	return nil
}

func ValidateToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	// cleanup empty header
	tokenString := strings.TrimSpace(header)
	if tokenString == "" {
		return "", errors.New("bearer token not found")
	}

	// separate bearer token, which comes in the form: "Bearer <token>"
	bearerToken := strings.Split(tokenString, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}

	return bearerToken[1], nil
}
