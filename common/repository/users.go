package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UsersRepository {
	return UsersRepository{db}
}

func (repo *UsersRepository) CreateUser(ctx context.Context, user model.User) error {
	if _, err := repo.db.ExecContext(ctx,
		`insert into users(id, role, username, password, fullname)
        values ($1, $2, $3, $4, $5)`,
		user.Id,
		user.Role,
		user.Username,
		user.Passhash,
		user.Fullname); err != nil {
		return err
	}

	return nil
}

func (repo *UsersRepository) FindUser(ctx context.Context, userId uuid.UUID) (model.User, error) {
	row := repo.db.QueryRowContext(ctx,
		`select id, role, username, password, fullname from users
        where id = $1`, userId)

	var user model.User
	if err := row.Scan(
		&user.Id,
		&user.Role,
		&user.Username,
		&user.Passhash,
		&user.Fullname); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (repo *UsersRepository) FindAllUsers(
	ctx context.Context,
	cursor uuid.UUID,
) ([]model.User, error) {
	var rows *sql.Rows
	var err error

	if (cursor != uuid.UUID{}) {
		rows, err = repo.db.QueryContext(ctx,
			`select id, role, username, password, fullname from users
            where id > $1
            order by id
            limit 10`, cursor)
	} else {
		rows, err = repo.db.QueryContext(ctx,
			`select id, role, username, password, fullname from users
            order by id
            limit 10`)
	}
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0, 10)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.Id,
			&user.Role,
			&user.Username,
			&user.Passhash,
			&user.Fullname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UsersRepository) UpdateUser(ctx context.Context, user model.User) error {
	if _, err := repo.db.ExecContext(ctx,
		`update users set
        username = coalesce(nullif($1, ''), username),
        fullname = coalesce(nullif($2, ''), fullname),
        password = coalesce(nullif($3, ''), password)
        where id = $4`, user.Username, user.Fullname, user.Passhash, user.Id); err != nil {
		return err
	}

	return nil
}
