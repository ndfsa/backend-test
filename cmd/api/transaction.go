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
		userId, _ := token.GetUserId(r)
		transactionIdString := r.URL.Query().Get("id")

		if len(transactionIdString) == 0 {
			err := repository.GetTransactions(userId)
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

		err = repository.GetTransaction(userId, transactionId)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		// return one transaction
	})
}

func executeTransaction(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, _ := token.GetUserId(r)

		var transaction dto.TransactionDto
		if err := util.Receive(r.Body, &transaction); err != nil {
			util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}

		if err := repository.ExecuteTransaction(r.Context(), db, userId, transaction); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

func rollbackTransaction(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := repository.RollbackTransaction(0, 0); err != nil {
			util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}
	})
}
