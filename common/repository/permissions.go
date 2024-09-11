package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type OwnershipRepository struct {
	db *sql.DB
}

var ErrOwnership = errors.New("user does not own resource")

func NewOwnershipRepository(db *sql.DB) OwnershipRepository {
	return OwnershipRepository{db}
}

func (repo *OwnershipRepository) CheckUserOwnership(
	userId, owner uuid.UUID,
) error {
	if userId == owner {
		return ErrOwnership
	}
	return nil
}

func (repo *OwnershipRepository) CheckServiceOwnership(
	ctx context.Context, serviceId, userId uuid.UUID,
) error {
	row := repo.db.QueryRowContext(ctx, `select exists(
        select 1 from user_service where (user_id, service_id) = ($1, $2))`,
		userId, serviceId)
	var owns bool
	if err := row.Scan(&owns); err != nil {
		return err
	}

	if !owns {
		return ErrOwnership
	}
	return nil
}

func (repo *OwnershipRepository) CheckTransactionOwnership(
	ctx context.Context, transactionId, userId uuid.UUID,
) error {
	row := repo.db.QueryRowContext(ctx, `select exists(
        select 1 from user_service us
        join transactions t
        on us.service_id = t.source
        or us.service_id = t.destination
        where us.user_id = $1
        and t.id = $2)`, userId, transactionId)
	var owns bool
	if err := row.Scan(&owns); err != nil {
		return err
	}

	if !owns {
		return ErrOwnership
	}
	return nil
}
