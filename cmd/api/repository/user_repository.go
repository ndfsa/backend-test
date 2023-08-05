package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/ndfsa/backend-test/cmd/api/dto"
	ilib "github.com/ndfsa/backend-test/internal"
)

func decodeJson(body io.ReadCloser) (dto.UserDto, error) {
	var userDto dto.UserDto

	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&userDto); err != nil {
		return userDto, errors.New("invalid parameters")
	}

	return userDto, nil
}

func CreateUser(db *sql.DB, body io.ReadCloser) (uint64, error) {
	newUser, err := decodeJson(body)
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

func ReadUser(db *sql.DB, userId uint64) (ilib.User, error) {
	var user ilib.User

	row := db.QueryRow("SELECT id, fullname, username FROM users WHERE id = $1", userId)
	if err := row.Err(); err != nil {
		return user, err
	}

	if err := row.Scan(&user.UserId, &user.Fullname, &user.Username); err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser(db *sql.DB, body io.ReadCloser, userId uint64) error {
	user, err := decodeJson(body)
	if err != nil {
		return err
	}

    row := db.QueryRow("CALL UPDATE_USER($1, $2, $3)", userId, user.Fullname, user.Username)
    if err := row.Err(); err != nil {
        return err
    }

	return nil
}

func DeleteUser(db *sql.DB, userId uint64) error {
    if userId == 1 {
        return errors.New("cannot delete root user")
    }

    row := db.QueryRow("DELETE FROM users WHERE id = $1", userId)
    if err := row.Err(); err != nil {
        return err
    }

	return nil
}
