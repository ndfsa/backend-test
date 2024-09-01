package token

import (
	"errors"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

const (
	USER_ID_KEY = "userId"
	REFRESH_KEY = "refreshKey"
    KEY = "2bbb515c1311dd69a609a0d553dc7ac1ac8eadc2b22daa9aaa99483d2f381374"
)

func GetUserId(bearerToken string) (uuid.UUID, error) {
    _, encodedToken, _ := strings.Cut(bearerToken, " ")
	parser := paseto.NewParserWithoutExpiryCheck()
	key, err := paseto.V4SymmetricKeyFromHex(KEY)
	if err != nil {
		return uuid.UUID{}, err
	}

	token, err := parser.ParseV4Local(key, encodedToken, nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	encodedId, err := token.GetString(USER_ID_KEY)
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(encodedId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func ValidateAccessToken(bearerToken string, hexKey string) (uuid.UUID, error) {
    _, encodedToken, found := strings.Cut(bearerToken, " ")
    if !found {
        return uuid.UUID{}, errors.New("invalid bearer token")
    }

	parser := paseto.NewParserForValidNow()

	key, err := paseto.V4SymmetricKeyFromHex(hexKey)
	if err != nil {
		return uuid.UUID{}, err
	}

	token, err := parser.ParseV4Local(key, encodedToken, nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	encodedId, err := token.GetString(USER_ID_KEY)
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(encodedId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func GenerateAccessToken(userId uuid.UUID) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(15 * time.Minute))

	token.SetString(USER_ID_KEY, userId.String())

	key, err := paseto.V4SymmetricKeyFromHex(KEY)
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(key, nil), nil
}

func GenerateRefreshToken(userId uuid.UUID) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(30 * 24 * time.Hour))

	token.SetString(USER_ID_KEY, userId.String())
	token.Set(REFRESH_KEY, true)

	key, err := paseto.V4SymmetricKeyFromHex(KEY)
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(key, nil), nil
}
