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

func NewSrvRepository(db *sql.DB) ServicesRepository {
	return ServicesRepository{db}
}

func (repo *ServicesRepository) CreateService(
	ctx context.Context, service model.Service,
) error {
	if _, err := repo.db.ExecContext(ctx,
		`insert into services(id, type, state, permissions, currency, init_balance, balance)
        values ($1, $2, $3, $4, $5, $6, $7)`,
		service.Id,
		service.Type,
		service.State,
		service.Permissions,
		service.Currency,
		service.InitBalance,
		service.Balance); err != nil {
		return err
	}

	return nil
}

func (repo *ServicesRepository) LinkServiceToUser(
	ctx context.Context, serviceId, userId uuid.UUID,
) error {
	if _, err := repo.db.ExecContext(ctx,
		`insert into user_service(user_id, service_id)
        values ($1, $2)`,
		userId,
		serviceId); err != nil {
		return err
	}

	return nil
}

func (repo *ServicesRepository) FindService(
	ctx context.Context, id uuid.UUID,
) (model.Service, error) {
	row := repo.db.QueryRowContext(ctx,
		`select id, type, state, permissions, currency, init_balance, balance
        from services where id = $1`, id)

	var service model.Service
	if err := row.Scan(
		&service.Id,
		&service.Type,
		&service.State,
		&service.Permissions,
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
			`select id, type, state, permissions, currency, init_balance, balance from services
            where id > $1
            order by id
            limit 10`, cursor)
	} else {
		rows, err = repo.db.QueryContext(ctx,
			`select id, type, state, permissions, currency, init_balance, balance from services
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
			&service.Permissions,
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

func (repo *ServicesRepository) UpdateService(
	ctx context.Context, service model.Service,
) error {
	result, err := repo.db.ExecContext(ctx,
		"update services set state = $1 where id = $2", service.State, service.Id)
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

func (repo *ServicesRepository) FindUserServices(
	ctx context.Context, user uuid.UUID,
) ([]model.Service, error) {
	rows, err := repo.db.QueryContext(ctx,
		`select s.id, s.type, s.state, s.permissions, s.currency, s.init_balance, s.balance
        from services s
        join user_service us on us.service_id = s.id
        where us.user_id = $1
        order by id`, user)
	if err != nil {
		return nil, err
	}

	services := make([]model.Service, 0)
	for rows.Next() {
		var service model.Service
		if err := rows.Scan(
			&service.Id,
			&service.Type,
			&service.State,
			&service.Permissions,
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
