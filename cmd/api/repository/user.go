package repository

import (
	"database/sql"

	"context"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/cmd/api/dto"
	"github.com/ndfsa/cardboard-bank/internal/model"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return UsersRepository{db}
}

func (r *UsersRepository) Get(ctx context.Context, userId uuid.UUID) (model.User, error) {
	var user model.User

	row := r.db.QueryRowContext(ctx,
		"SELECT id, fullname, username FROM users WHERE id = $1", userId)
	if err := row.Err(); err != nil {
		return user, err
	}

	if err := row.Scan(&user.UserId, &user.Fullname, &user.Username); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UsersRepository) Update(ctx context.Context, user dto.UserDto, userId uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE users SET fullname = $1, username = $2 WHERE id = $3",
		user.Fullname, user.Username, userId); err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) Delete(ctx context.Context, userId uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId); err != nil {
		return err
	}

	return nil
}
