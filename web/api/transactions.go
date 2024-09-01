package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
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
		var req dto.CreateTransactionRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		transaction, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		if err := factory.repo.CreateTransaction(r.Context(), transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(dto.CreateTransactionResponseDTO{
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

		transaction, err := factory.repo.FindTransaction(r.Context(), transactionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(dto.ReadTransactionResponseDTO{
			Id:          transaction.Id.String(),
			State:       transaction.State,
			Time:        transaction.Time,
			Currency:    transaction.Currency,
			Amount:      transaction.Amount.String(),
			Source:      transaction.Source.String(),
			Destination: transaction.Destination.String(),
		}); err != nil {
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

		transactions, err := factory.repo.FindAllTransactions(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		res := make([]dto.ReadTransactionResponseDTO, 0, len(transactions))
		for _, transaction := range transactions {
			res = append(res, dto.ReadTransactionResponseDTO{
				Id:          transaction.Id.String(),
				State:       transaction.State,
				Time:        transaction.Time,
				Currency:    transaction.Currency,
				Amount:      transaction.Amount.String(),
				Source:      transaction.Source.String(),
				Destination: transaction.Destination.String(),
			})
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}
