package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/ndfsa/backend-test/cmd/api/dto"
	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)


func CreateTransactionRoutes(db *sql.DB, baseUrl string) {
    repo := repository.NewTransactionsRepository(db)

	http.Handle(baseUrl+"/transaction", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getTransaction(repo)))
	http.Handle(baseUrl+"/transaction/execute", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPost),
		middleware.Auth)(executeTransaction(repo)))
	http.Handle(baseUrl+"/transaction/rollback", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(rollbackTransaction(repo)))
}

func getTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, _ := token.GetUserId(r)
		transactionIdString := r.URL.Query().Get("id")

		if len(transactionIdString) == 0 {
			err := repo.GetAll(userId)
			if err != nil {
				util.Error(&w, http.StatusInternalServerError, err.Error())
				return
			}
			// return transaction list
		}

		transactionId, err := strconv.ParseUint(transactionIdString, 10, 64)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		err = repo.Get(userId, transactionId)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		// return one transaction
	})
}

func executeTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, _ := token.GetUserId(r)

		var transaction dto.TransactionDto
		if err := util.Receive(r.Body, &transaction); err != nil {
			util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}

		if err := repo.Execute(r.Context(), userId, transaction); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

func rollbackTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := repo.Rollback(0, 0); err != nil {
			util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}
	})
}
