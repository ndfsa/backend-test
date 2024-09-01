package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
)

type TransactionsRepository struct {
	db      *sql.DB
	pool    chan model.Transaction
	workers int
	queue   int
}

func (repo *TransactionsRepository) worker(jobs <-chan model.Transaction) {
	for transaction := range jobs {
		if err := repo.executeTransaction(transaction); err != nil {
			continue
		}
	}
}

func NewTransactionsRepository(db *sql.DB, queue, workers int) TransactionsRepository {
	if queue < 0 {
		panic("invalid queue size")
	}
	if workers < 1 {
		panic("number of workers must be >= 1")
	}

	pool := make(chan model.Transaction, queue)

	repo := TransactionsRepository{db, pool, workers, queue}

	for i := 0; i < workers; i++ {
		go repo.worker(pool)
	}

	return repo
}

func (repo *TransactionsRepository) RegisterTransaction(
	ctx context.Context,
	transaction model.Transaction,
) error {
	row := repo.db.QueryRowContext(ctx,
		`insert into transactions(id, state, time, currency, amount, source, destination)
        values($1, $2, $3, $4, $5, $6, $7) returning time`,
		transaction.Id,
		transaction.State,
		transaction.Time,
		transaction.Currency,
		transaction.Amount,
		transaction.Source,
		transaction.Destination)

	if err := row.Scan(&transaction.Time); err != nil {
		return err
	}

	select {
	case repo.pool <- transaction:
		return nil
	default:
		return errors.New("pool full")
	}
}

func (repo *TransactionsRepository) executeTransaction(
	transaction model.Transaction,
) error {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`select balance from services
        where id = $1 or id = $2 for no key update`,
		transaction.Source,
		transaction.Destination); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`update services set
        balance = balance - $1
        where id = $2`,
		transaction.Amount,
		transaction.Source); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`update services set
        balance = balance + $1
        where id = $2`,
		transaction.Amount,
		transaction.Destination); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`update transactions set
        state = $1
        where id = $2`, model.TransactionStateSuccess, transaction.Id); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (repo *TransactionsRepository) GetTransaction(
	ctx context.Context,
	id uuid.UUID,
) (model.Transaction, error) {
	row := repo.db.QueryRowContext(ctx,
		`select id, state, time, currency, amount, source, destination
        from transactions where id = $1`, id)

	var transaction model.Transaction
	if err := row.Scan(
		&transaction.Id,
		&transaction.State,
		&transaction.Time,
		&transaction.Currency,
		&transaction.Amount,
		&transaction.Source,
		&transaction.Destination); err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func (repo *TransactionsRepository) GetTransactions(
	ctx context.Context,
	cursor uuid.UUID,
) ([]model.Transaction, error) {
	var rows *sql.Rows
	var err error

	if (cursor != uuid.UUID{}) {
		rows, err = repo.db.QueryContext(ctx,
			`select id, state, time, currency, amount, source, destination from transactions
            where id > $1
            order by id
            limit 10`, cursor)
	} else {
		rows, err = repo.db.QueryContext(ctx,
			`select id, state, time, currency, amount, source, destination
            from transactions
            order by id
            limit 10`)
	}
	if err != nil {
		return nil, err
	}

	transactions := make([]model.Transaction, 0, 10)
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(
			&transaction.Id,
			&transaction.State,
			&transaction.Time,
			&transaction.Currency,
			&transaction.Amount,
			&transaction.Source,
			&transaction.Destination); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
