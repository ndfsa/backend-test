package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type TransactionsHandlerFactory struct {
	repo repository.TransactionsRepository
	mdf  middleware.MiddlewareFactory
}

func NewTransactionsHandlerFactory(
	repo repository.TransactionsRepository,
	mdf middleware.MiddlewareFactory,
) TransactionsHandlerFactory {
	return TransactionsHandlerFactory{repo, mdf}
}

func (factory *TransactionsHandlerFactory) CreateTransaction() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateTransactionRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		transaction, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		if err := factory.repo.RegisterTransaction(r.Context(), transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(CreateTransactionResponseDTO{
			Id: transaction.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *TransactionsHandlerFactory) ReadSingleTransaction() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		transactionId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		transaction, err := factory.repo.GetTransaction(r.Context(), transactionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *TransactionsHandlerFactory) ReadMultipleTransactions() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return err
			}
		} else {
			cursor = uuid.UUID{}
		}

		transactions, err := factory.repo.GetTransactions(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}
