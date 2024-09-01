package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
)

type ServicesRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) ServicesRepository {
	return ServicesRepository{db}
}

func (repo *ServicesRepository) CreateService(
	ctx context.Context, service model.Service,
) error {
	if _, err := repo.db.ExecContext(ctx,
		`insert into services(id, type, state, currency, init_balance, balance)
        values ($1, $2, $3, $4, $5, $6)`,
		service.Id,
		service.Type,
		service.State,
		service.Currency,
		service.InitBalance,
		service.Balance); err != nil {
		return err
	}

	return nil
}

func (repo *ServicesRepository) FindService(
	ctx context.Context, id uuid.UUID,
) (model.Service, error) {
	row := repo.db.QueryRowContext(ctx,
		"select id, type, state, currency, init_balance, balance from services where id = $1",
		id)

	var service model.Service
	if err := row.Scan(
		&service.Id,
		&service.Type,
		&service.State,
		&service.Currency,
		&service.InitBalance,
		&service.Balance); err != nil {
		return model.Service{}, err
	}

	return service, nil
}

func (repo *ServicesRepository) FindAllServices(
	ctx context.Context, cursor uuid.UUID,
) ([]model.Service, error) {
	var rows *sql.Rows
	var err error

	if (cursor != uuid.UUID{}) {
		rows, err = repo.db.QueryContext(ctx,
			`select id, type, state, currency, init_balance, balance from services
            where id > $1
            order by id
            limit 10`, cursor)
	} else {
		rows, err = repo.db.QueryContext(ctx,
			`select id, type, state, currency, init_balance, balance from services
            order by id
            limit 10`)
	}
	if err != nil {
		return nil, err
	}

	services := make([]model.Service, 0, 10)
	for rows.Next() {
		var service model.Service
		if err := rows.Scan(
			&service.Id,
			&service.Type,
			&service.State,
			&service.Currency,
			&service.InitBalance,
			&service.Balance); err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (repo *ServicesRepository) SetServiceState(
	ctx context.Context, id uuid.UUID, state string,
) error {
	result, err := repo.db.ExecContext(ctx,
		"update services set state = $1 where id = $2", state, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("%d rows changed", rows)
	}
	return nil
}
