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
