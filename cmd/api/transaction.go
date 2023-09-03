package main

import (
	"database/sql"
	"net/http"

	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
)

func CreateTransactionRoutes(db *sql.DB, baseUrl string) {
	http.Handle(baseUrl+"/transaction", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getTransaction(db)))
	http.Handle(baseUrl+"/transaction/execute", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPost),
		middleware.Auth)(executeTransaction(db)))
	http.Handle(baseUrl+"/transaction/rollback", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(rollbackTransaction(db)))
}

func getTransaction(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        repository.GetTransaction()
	})
}

func executeTransaction(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        repository.ExecuteTransaction()
	})
}

func rollbackTransaction(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        repository.RollbackTransaction()
	})
}
