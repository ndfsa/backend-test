package repository

import "github.com/shopspring/decimal"

func ExecuteTransaction(
	userId uint64,
	from uint64,
	to uint64,
	amount decimal.Decimal,
	currency string) error {
	return nil
}

func GetTransaction() error {
	return nil
}

func GetTransactions() error {
	return nil
}

func RollbackTransaction() {

}
