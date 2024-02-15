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

func CreateServiceRoutes(db *sql.DB, baseUrl string, key string) {
	repo := repository.NewServicesRepository(db)

	http.Handle("GET "+baseUrl+"/service", middleware.Chain(
		middleware.Logger,
		middleware.Auth(key))(get(repo)))
	http.Handle("POST "+baseUrl+"/service/create", middleware.Chain(
		middleware.Logger,
		middleware.Auth(key))(create(repo)))
	http.Handle("DELETE "+baseUrl+"/service/cancel", middleware.Chain(
		middleware.Logger,
		middleware.Auth(key))(cancel(repo)))
}

func get(repo repository.ServicesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user id from header
		userId, _ := token.GetUserId(r)

		// get service param from query
		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			// if no service ID is found in query, return all services
			services, err := repo.GetAll(r.Context(), userId)
			if err != nil {
				util.SendError(&w, http.StatusInternalServerError, err.Error())
				return
			}

			util.Send(&w, services)
			return
		}
		serviceId, err := strconv.ParseUint(serviceIdQuery, 10, 64)
		if err != nil {
			util.SendError(&w, http.StatusBadRequest, err.Error())
			return
		}

		// get all user services
		services, err := repo.Get(r.Context(), userId, serviceId)
		if err != nil {
			util.SendError(&w, http.StatusNoContent, err.Error())
			return
		}

		util.Send(&w, services)
	})
}

func create(repo repository.ServicesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user ID
		userId, _ := token.GetUserId(r)

		// get service details from request
		var service dto.ServiceDto
		util.Receive[dto.ServiceDto](r.Body, &service)

		// call repository
		serviceId, err := repo.Create(r.Context(), userId, service)
		if err != nil {
			util.SendError(&w, http.StatusInternalServerError, err.Error())
		}

		// return service ID
		util.Send(&w, serviceId)
	})
}

func cancel(repo repository.ServicesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID
		userId, _ := token.GetUserId(r)
		// get service param from query
		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			// if no service ID is found in query, return all services
			util.SendError(&w, http.StatusBadRequest, "no service id")
			return
		}
		serviceId, err := strconv.ParseUint(serviceIdQuery, 10, 64)
		if err != nil {
			util.SendError(&w, http.StatusBadRequest, err.Error())
			return
		}
		if err := repo.Cancel(r.Context(), userId, serviceId); err != nil {
			util.SendError(&w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
