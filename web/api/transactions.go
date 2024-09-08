package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type TransactionsHandlerFactory struct {
	repo repository.TransactionsRepository
}

func NewTransactionsHandlerFactory(
	repo repository.TransactionsRepository,
) TransactionsHandlerFactory {
	return TransactionsHandlerFactory{repo}
}

func (factory *TransactionsHandlerFactory) CreateTransaction() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTransactionRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		transaction, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.CreateTransaction(r.Context(), transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.CreateTransactionResponseDTO{
			Id: transaction.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *TransactionsHandlerFactory) ReadSingleTransaction() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		transactionId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		transaction, err := factory.repo.FindTransaction(r.Context(), transactionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
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
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *TransactionsHandlerFactory) ReadMultipleTransactions() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				return
			}
		} else {
			cursor = uuid.UUID{}
		}

		transactions, err := factory.repo.FindAllTransactions(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
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
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *TransactionsHandlerFactory) ReadServiceTransactions() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}

		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				return
			}
		} else {
			cursor = uuid.UUID{}
		}

		transactions, err := factory.repo.FindServiceTransactions(r.Context(), serviceId, cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
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
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}
