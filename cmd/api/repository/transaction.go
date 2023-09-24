package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ndfsa/backend-test/cmd/api/dto"
)

func ExecuteTransaction(
	ctx context.Context,
	db *sql.DB,
	userId uint64,
	transaction dto.TransactionDto) error {

    // check if user owns source service
    row := db.QueryRowContext(ctx, `SELECT EXISTS (
        SELECT 1 FROM users u
        JOIN user_service us ON u.id = us.user_id
        JOIN services s ON s.id = us.service_id
        WHERE u.id = $1 AND s.id = $2)`, userId, transaction.From)
    if err := row.Err(); err != nil {
        return err
    }
    var belongsToUser bool
    if err := row.Scan(&belongsToUser); err != nil {
        return err
    }

    if !belongsToUser {
        return errors.New("source service does not belong to user")
    }

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// create the transaction
	row = tx.QueryRowContext(ctx, `INSERT INTO transactions(state, currency, amount)
        VALUES('DONE', $2, $1)
        RETURNING id`, transaction.Amount, transaction.Currency)
	if err := row.Err(); err != nil {
		return err
	}
	var transactionId uint64
	if err := row.Scan(&transactionId); err != nil {
		return err
	}

	// create the relation between the transaction and the corresponding services
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO service_transaction(transaction_id, from_service_id, to_service_id, user_id),
        VALUES($1, $2, $3, $4)`,
		transactionId, transaction.From, transaction.To, userId); err != nil {
		return err
	}

	// update service balances
	if _, err := tx.ExecContext(ctx, `UPDATE services AS s
        SET balance = balance + c.amount
        FROM (values ($2, -$1), ($3, $1))
        AS c(srv_id, amount)
        WHERE c.srv_id = s.id`,
		transaction.Amount, transaction.From, transaction.To, userId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetTransaction(userId uint64, transactionId uint64) error {
	return nil
}

func GetTransactions(userId uint64) error {
	return nil
}

func RollbackTransaction(userId uint64, transactionId uint64) error {
	return nil
}
