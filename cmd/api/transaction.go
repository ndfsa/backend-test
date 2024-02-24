package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/cmd/api/dto"
	"github.com/ndfsa/cardboard-bank/cmd/api/repository"
	"github.com/ndfsa/cardboard-bank/internal/encoding"
	"github.com/ndfsa/cardboard-bank/internal/token"
)

func getTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
		}

		transactionIdString := r.PathValue("id")
		transactionId, err := uuid.Parse(transactionIdString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		err = repo.Get(userId, transactionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}

func getAllTransactions(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
		}

		transactions, err := repo.GetAll(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

        encoding.Send(w, transactions)
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
			log.Println(err)
			return
		}

		if err := repo.Execute(r.Context(), userId, transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}

func rollbackTransaction(repo repository.TransactionsRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        encodedToken := r.Header.Get("Authorization")

        userId, err := token.GetUserId(encodedToken, tokenKey)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            log.Println(err)
            return
        }

        transactionIdString := r.PathValue("id")
        transactionId, err := uuid.Parse(transactionIdString)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            log.Println(err)
            return
        }

		if err := repo.Rollback(userId, transactionId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
	})
}
