package repository

import (
	"database/sql"
	"errors"
	"io"

	"github.com/ndfsa/backend-test/cmd/api/dto"
	"github.com/ndfsa/backend-test/internal/model"
	"github.com/ndfsa/backend-test/internal/util"
)

func ReadUser(db *sql.DB, userId uint64) (model.User, error) {
	var user model.User

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
	var user dto.UserDto
	err := util.Receive(body, &user)
	if err != nil {
		return err
	}

	if _, err := db.Exec("UPDATE users SET fullname = $1, username = $2 WHERE id = $3",
		user.Fullname, user.Username, userId); err != nil {
		return err
	}

	return nil
}

func DeleteUser(db *sql.DB, userId uint64) error {
	if userId == 1 {
		return errors.New("cannot delete root user")
	}

	if _, err := db.Exec("DELETE FROM users WHERE id = $1", userId); err != nil {
		return err
    }

	return nil
}
