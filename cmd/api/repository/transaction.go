package repository

import (
	"github.com/ndfsa/backend-test/cmd/api/dto"
)

func ExecuteTransaction(userId uint64, transaction dto.TransactionDto) error {
	return nil
}

func GetTransaction() error {
	return nil
}

func GetTransactions() error {
	return nil
}

func RollbackTransaction(userId uint64, transactionId uint64) error {
	return nil
}
