package repository

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"io"

	"github.com/ndfsa/backend-test/cmd/auth/dto"
	"github.com/ndfsa/backend-test/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(db *sql.DB, username string, password string) (uint64, error) {
	row := db.QueryRow("SELECT password, id FROM users WHERE username = $1", username)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var storedPassword string
	var id uint64
	if err := row.Scan(&storedPassword, &id); err != nil {
		return 0, err
	}

    passwordBytes := []byte(password)
	reducedPassword := make([]byte, base64.StdEncoding.EncodedLen(len(passwordBytes)))
    base64.StdEncoding.Encode(reducedPassword, passwordBytes)
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), reducedPassword); err != nil {
		return 0, err
	}

	return id, nil
}

func SignUp(db *sql.DB, body io.ReadCloser) (uint64, error) {
	var newUser dto.SignUpDTO
	if err := util.Receive(body, &newUser); err != nil {
		return 0, err
	}

	if newUser.Username == "" || newUser.Password == "" || newUser.Fullname == "" {
		return 0, errors.New("Invalid data")
	}

	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return 0, err
	}

	row := db.QueryRow(
		"INSERT INTO users(fullname, username, password) VALUES ($1, $2, $3) RETURNING id",
		newUser.Fullname, newUser.Username, hashedPassword)

	if err := row.Err(); err != nil {
		return 0, err
	}

	var res uint64
	if err := row.Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}

func hashPassword(password string) (string, error) {
	reducedString := base64.StdEncoding.EncodeToString([]byte(password))
	hashedString, err := bcrypt.GenerateFromPassword([]byte(reducedString), 0)
	if err != nil {
		return "", err
	}
	return string(hashedString), nil
}
