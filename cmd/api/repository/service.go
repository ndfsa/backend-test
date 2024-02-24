package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/internal/model"
)

type ServicesRepository struct {
	db *sql.DB
}

func NewServicesRepository(db *sql.DB) ServicesRepository {
	return ServicesRepository{db}
}

func (r *ServicesRepository) GetServiceTransactions(
	ctx context.Context,
	userId uuid.UUID,
	serviceId uuid.UUID) ([]model.Transaction, error) {

	row := r.db.QueryRowContext(ctx,
		`SELECT true FROM user_service
        WHERE (userId, serviceId) = ($1, $2)`, userId, serviceId)
	if err := row.Err(); err != nil {
		return nil, err
	}

    var exists bool
    if err := row.Scan(&exists); err != nil {
        return nil, err
    }
    if !exists {
        return nil, errors.New("user does not have service")
    }

	rows, err := r.db.QueryContext(ctx,
		`SELECT t.id, t.state, t.currency, t.amount, t.from, t.to
        FROM transactions t
        WHERE t.from = $1
        OR t.to = $1`, serviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]model.Transaction, 0)
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(
			&transaction.Id,
			&transaction.State,
			&transaction.Currency,
			&transaction.Amount,
			&transaction.From,
			&transaction.To); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *ServicesRepository) GetAll(
	ctx context.Context,
	userId uuid.UUID) ([]model.Service, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.type, s.state, s.currency, s.init_balance, s.balance
        FROM users u
        JOIN user_service us ON u.id = us.user_id
        JOIN services s ON s.id = us.service_id
        WHERE u.id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]model.Service, 0)
	for rows.Next() {
		var service model.Service
		if err := rows.Scan(
			&service.Id,
			&service.Type,
			&service.State,
			&service.Currency,
			&service.InitBalance,
			&service.Balance); err != nil {
			return services, err
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return services, err
	}

	return services, nil
}

func (r *ServicesRepository) Get(
	ctx context.Context,
	userId uuid.UUID,
	serviceId uuid.UUID) (model.Service, error) {

	rows := r.db.QueryRowContext(ctx,
		`SELECT s.id, s.type, s.state, s.currency, s.init_balance, s.balance
        FROM users u
        JOIN user_service us ON u.id = us.user_id
        JOIN services s ON s.id = us.service_id
        WHERE u.id = $1 AND s.id = $2`, userId, serviceId)

	var service model.Service

	if err := rows.Err(); err != nil {
		return service, err
	}

	if err := rows.Scan(
		&service.Id,
		&service.Type,
		&service.State,
		&service.Currency,
		&service.InitBalance,
		&service.Balance); err != nil {
		return service, err
	}

	return service, nil
}

func (r *ServicesRepository) Create(
	ctx context.Context,
	userId uuid.UUID,
	service model.Service) (uuid.UUID, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer tx.Rollback()

	serviceId, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO services (id, type, state, currency, init_balance, balance)
        VALUES ($1, $2, 'REQ', $3, $4, $5)`,
		serviceId,
		service.Type,
		service.Currency,
		0,
		service.Balance); err != nil {
		return uuid.UUID{}, err
	}

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO user_service (user_id, service_id) VALUES ($1, $2)",
		userId,
		serviceId); err != nil {
		return uuid.UUID{}, err
	}

	if err = tx.Commit(); err != nil {
		return uuid.UUID{}, err
	}

	return serviceId, nil
}

func (r *ServicesRepository) Cancel(
	ctx context.Context,
	userId uuid.UUID,
	serviceId uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx,
		`UPDATE services SET state = 'CLD'
        FROM users JOIN user_service ON users.id = user_id
        WHERE users.id = $1 AND services.id = $2`,
		userId,
		serviceId); err != nil {
		return err
	}

	return nil
}
