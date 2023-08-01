package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/ndfsa/backend-test/cmd/back/dto"
	"github.com/ndfsa/backend-test/internal/models"
)

func CreateUser(db *sql.DB, body io.ReadCloser) (uint64, error) {
	var newUser dto.UserDto

	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&newUser); err != nil ||
		newUser.Name == "" ||
		newUser.Password == "" ||
		newUser.Username == "" {

		return 0, errors.New("invalid parameters")
	}

	row := db.QueryRow("select selectUser($1, $2, $3)",
		newUser.Name,
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

func ReadUser(db *sql.DB, userId uint64) (models.User, error) {
	var user models.User

	row, err := db.Query(
		"SELECT userId, userFullName, userName FROM album WHERE userId = ?", userId)
	if err != nil {
		return user, err
	}

	if err := row.Scan(&user.UserId, &user.Fullname, &user.Username); err != nil {
		return user, err
	}

	return user, nil
}

func UpdateUser() error {
	return errors.New("")
}

func DeleteUser() {
}
