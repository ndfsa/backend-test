package token

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/ndfsa/cardboard-bank/common/model"
)

const (
	USER_KEY    = "userKey"
	REFRESH_KEY = "refreshKey"
	KEY         = "2bbb515c1311dd69a609a0d553dc7ac1ac8eadc2b22daa9aaa99483d2f381374"
)

func ValidateAccessToken(bearerToken string) (model.User, error) {
	_, encodedToken, found := strings.Cut(bearerToken, " ")
	if !found {
		return model.User{}, errors.New("invalid bearer token")
	}

	parser := paseto.NewParserForValidNow()

	key, err := paseto.V4SymmetricKeyFromHex(KEY)
	if err != nil {
		return model.User{}, err
	}

	token, err := parser.ParseV4Local(key, encodedToken, nil)
	if err != nil {
		return model.User{}, err
	}

	encodedUser, err := token.GetString(USER_KEY)
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(encodedUser), &user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func GenerateAccessToken(user model.User) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(15 * time.Minute))

	payload, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	token.SetString(USER_KEY, string(payload))

	key, err := paseto.V4SymmetricKeyFromHex(KEY)
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(key, nil), nil
}

func GenerateRefreshToken(user model.User) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(30 * 24 * time.Hour))

	payload, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	token.SetString(USER_KEY, string(payload))
	token.Set(REFRESH_KEY, true)

	key, err := paseto.V4SymmetricKeyFromHex(KEY)
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(key, nil), nil
}
