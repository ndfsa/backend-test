package repository

import (
	"database/sql"
	"io"

	"github.com/ndfsa/backend-test/cmd/auth/dto"
	"github.com/ndfsa/backend-test/internal/util"
)

func AuthenticateUser(db *sql.DB, username string, password string) (uint64, error) {
	row := db.QueryRow("SELECT AUTHENTICATE_USER($1, $2)", username, password)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func SignUp(db *sql.DB, body io.ReadCloser) (uint64, error) {
	var newUser dto.SignUpDTO
	if err := util.Receive(body, &newUser); err != nil {
		return 0, err
	}

	row := db.QueryRow("SELECT CREATE_USER($1, $2, $3)",
		newUser.Fullname,
		newUser.Username,
		newUser.Password)

	if err := row.Err(); err != nil {
		return 0, err
	}

	var res uint64
	if err := row.Scan(&res); err != nil {
		return 0, err
	}

	return res, nil
}
