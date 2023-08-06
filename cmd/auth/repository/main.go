package repository

import (
	"database/sql"
	"errors"
	"io"
	"log"

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
	err := util.DecodeJson(body, &newUser)
	if err != nil {
		return 0, err
	}

	row := db.QueryRow("SELECT CREATE_USER($1, $2, $3)",
		newUser.Fullname,
		newUser.Username,
		newUser.Password)

	if err := row.Err(); err != nil {
		log.Println(err.Error())
		return 0, errors.New("could not create user")
	}

	var res uint64
	if err := row.Scan(&res); err != nil {
		log.Println(err.Error())
		return 0, errors.New("could not create user")
	}

	return res, nil
}
