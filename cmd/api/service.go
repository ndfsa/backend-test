package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/cmd/api/dto"
	"github.com/ndfsa/cardboard-bank/cmd/api/repository"
	"github.com/ndfsa/cardboard-bank/internal/encoding"
	"github.com/ndfsa/cardboard-bank/internal/model"
	"github.com/ndfsa/cardboard-bank/internal/token"
)

func getAll(repo repository.ServicesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		services, err := repo.GetAll(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		encoding.Send(w, services)
	})
}

func get(repo repository.ServicesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		serviceIdQuery := r.PathValue("id")
		serviceId, err := uuid.Parse(serviceIdQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		service, err := repo.Get(r.Context(), userId, serviceId)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			log.Println(err)
			return
		}

		encoding.Send(w, service)
	})
}

func create(repo repository.ServicesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		var serviceDto dto.CreateServiceRequest
        if err := encoding.Receive[dto.CreateServiceRequest](r, &serviceDto); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            log.Println(err)
            return
        }

		serviceId, err := repo.Create(r.Context(), userId, model.Service{
			Type:        serviceDto.Type,
			Currency:    serviceDto.Currency,
			Balance: serviceDto.InitBalance,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		encoding.Send(w, dto.CreateServiceResponse{
            Id: serviceId.String(),
        })
	})
}

func cancel(repo repository.ServicesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		serviceIdQuery := r.PathValue("id")
		serviceId, err := uuid.Parse(serviceIdQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := repo.Cancel(r.Context(), userId, serviceId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}
