package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/cmd/auth/dto"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return AuthRepository{db}
}

func (r *AuthRepository) AuthenticateUser(
	ctx context.Context,
	username string,
	password string) (uuid.UUID, error) {

	row := r.db.QueryRowContext(
		ctx,
		"SELECT password, id FROM users WHERE username = $1",
		username)
	if row.Err() != nil {
		return uuid.UUID{}, row.Err()
	}

	var storedPass string
	var id uuid.UUID
	if err := row.Scan(&storedPass, &id); err != nil {
		return uuid.UUID{}, err
	}

	passBytes := []byte(password)

	// reduce string to base 64 to avoid dealing with unicode
	reducedPass := make([]byte, base64.StdEncoding.EncodedLen(len(passBytes)))
	base64.StdEncoding.Encode(reducedPass, passBytes)
	if err := bcrypt.CompareHashAndPassword([]byte(storedPass), reducedPass); err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, newUser dto.SignUpDTO) error {
	if newUser.Username == "" || newUser.Password == "" || newUser.Fullname == "" {
		return errors.New("Invalid data")
	}

	// reduce string to base 64 to avoid dealing with unicode
	reducedPass := base64.StdEncoding.EncodeToString([]byte(newUser.Password))
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(reducedPass), 0)
	if err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx,
		"INSERT INTO users(id, fullname, username, password) VALUES ($1, $2, $3, $4)",
		uuid.New(), newUser.Fullname, newUser.Username, string(hashedPass)); err != nil {
		return err
	}

	return nil
}
