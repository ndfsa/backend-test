package repository

import (
	"context"
	"database/sql"

	"github.com/ndfsa/backend-test/cmd/api/dto"
)

func ExecuteTransaction(
	ctx context.Context,
	db *sql.DB,
	userId uint64,
	transaction dto.TransactionDto) error {

	if _, err := db.ExecContext(ctx, `update services as s
        set balance = balance + c.amount
        from (values ($2, -$1), ($3, $1))
        as c(to_id, amount)
        where c.to_id = s.id`, transaction.Amount, transaction.From, transaction.To); err != nil {
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
