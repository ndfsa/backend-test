package token

import (
	"net/http"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

const (
	USER_ID_KEY      = "userId"
	ACCESS_DURATION  = 15 * time.Minute
	REFRESH_DURATION = 30 * 24 * time.Hour
)

func GetUserId(encodedToken string, hexKey string) (uuid.UUID, error) {
	parser := paseto.NewParserWithoutExpiryCheck()
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

func ValidateAccessToken(encodedToken string, hexKey string) error {
	parser := paseto.NewParserForValidNow()

	key, err := paseto.V4SymmetricKeyFromHex(hexKey)
	if err != nil {
		return err
	}

	_, err = parser.ParseV4Local(key, encodedToken, nil)
	if err != nil {
		return err
	}
	return nil
}

func ValidateRefreshToken(encodedToken, hexKey string) error {
	return nil
}

func generateToken(userId uuid.UUID, tokenKey string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	token.SetString(USER_ID_KEY, userId.String())

	key, err := paseto.V4SymmetricKeyFromHex(tokenKey)
	if err != nil {
		return "", err
	}

	return token.V4Encrypt(key, nil), nil
}

func GenerateAccessToken(userId uuid.UUID, tokenKey string) (string, error) {
	return generateToken(userId, tokenKey, ACCESS_DURATION)
}

func GenerateRefreshToken(userId uuid.UUID, tokenKey string) (string, error) {
	return generateToken(userId, tokenKey, REFRESH_DURATION)
}

func ValidateToken(r *http.Request) error {
	return nil
}
