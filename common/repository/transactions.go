package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
)

type TransactionsRepository struct {
	db       *sql.DB
	jobQueue chan<- model.Transaction
}

func NewTrsRepository(
	db *sql.DB,
	jobQueue chan<- model.Transaction,
) TransactionsRepository {
	repo := TransactionsRepository{db, jobQueue}
	return repo
}

func (repo *TransactionsRepository) CreateTransaction(
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
	case repo.jobQueue <- transaction:
		return nil
	default:
		return errors.New("queue is full or unavailable")
	}
}

func (repo *TransactionsRepository) ExecuteTransaction(
	transaction model.Transaction,
) error {
	if transaction.Source == transaction.Destination {
		return fmt.Errorf("transaction: %s invalid, src and dst are the same", transaction.Id)
	}

	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	srcRow := tx.QueryRow(
		`select id, type, state, currency, init_balance, balance from services
        where id = $1 for no key update`,
		transaction.Source)

	var srcService model.Service
	if err := srcRow.Scan(
		&srcService.Id,
		&srcService.Type,
		&srcService.State,
		&srcService.Currency,
		&srcService.InitBalance,
		&srcService.Balance,
	); err != nil {
		return err
	}

	if err := srcService.Debit(transaction.Amount); err != nil {
		return err
	}

	dstRow := tx.QueryRow(
		`select id, type, state, currency, init_balance, balance from services
        where id = $1 for no key update`,
		transaction.Source)

	var dstService model.Service
	if err := dstRow.Scan(
		&dstService.Id,
		&dstService.Type,
		&dstService.State,
		&dstService.Currency,
		&dstService.InitBalance,
		&dstService.Balance,
	); err != nil {
		return err
	}

	if err := dstService.Credit(transaction.Amount); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`update services set balance = $1
        where id = $2`,
		srcService.Balance,
		srcService.Id); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`update services set balance = $1
        where id = $2`,
		dstService.Balance,
		dstService.Id); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`update transactions set state = $1
        where id = $2`, model.TransactionStateSuccess, transaction.Id); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (repo *TransactionsRepository) FindTransaction(
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

func (repo *TransactionsRepository) FindAllTransactions(
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
