package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ndfsa/cardboard-bank/cmd/api/dto"
	"github.com/ndfsa/cardboard-bank/cmd/api/repository"
	"github.com/ndfsa/cardboard-bank/internal/encoding"
	"github.com/ndfsa/cardboard-bank/internal/token"
)

func get(repo repository.ServicesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			services, err := repo.GetAll(r.Context(), userId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			encoding.Send(w, services)
			return
		}
		serviceId, err := strconv.ParseUint(serviceIdQuery, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		// get all user services
		services, err := repo.Get(r.Context(), userId, serviceId)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			log.Println(err)
			return
		}

		encoding.Send(w, services)
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

		var service dto.ServiceDto
		encoding.Receive[dto.ServiceDto](r, &service)

		serviceId, err := repo.Create(r.Context(), userId, service)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		encoding.Send(w, serviceId)
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

		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("no service id")
			return
		}

		serviceId, err := strconv.ParseUint(serviceIdQuery, 10, 64)
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
