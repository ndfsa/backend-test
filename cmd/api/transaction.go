package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/cmd/api/dto"
	"github.com/ndfsa/cardboard-bank/cmd/api/repository"
	"github.com/ndfsa/cardboard-bank/internal/encoding"
	"github.com/ndfsa/cardboard-bank/internal/token"
)

func CreateTransactionRoutes(db *sql.DB, baseUrl string, tokenKey string) {
}

func getTransactionById(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
		}

		transactionIdString := r.PathValue("transactionId")
		transactionId, err := uuid.Parse(transactionIdString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		err = repo.Get(userId, transactionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	})
}

func getTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		_, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
		}

		// transactions, err := repo.GetAll(userId)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Println(err.Error())
		// 	return
		// }
		//
	})
}

func executeTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		var transaction dto.ExecuteTransactionRequest
		if err := encoding.Receive(r, &transaction); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		if err := repo.Execute(r.Context(), userId, transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	})
}

func rollbackTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := repo.Rollback(uuid.UUID{}, uuid.UUID{}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
	})
}
