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

func CreateServiceRoutes(db *sql.DB, baseUrl string) {
	http.Handle(baseUrl+"/service", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getService(db)))
	http.Handle(baseUrl+"/service/create", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPost),
		middleware.Auth)(createService(db)))
	http.Handle(baseUrl+"/service/cancel", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(cancelService(db)))
}

func getService(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user id from header
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		// get service param from query
		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			// if no service ID is found in query, return all services
			services, err := repository.GetServices(db, userId)
			if err != nil {
				util.Error(&w, http.StatusInternalServerError, err.Error())
				return
			}

			util.Send(&w, services)
			return
		}
		serviceId, err := strconv.ParseUint(serviceIdQuery, 10, 64)
		if err != nil {
			util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}

		// get all user services
		services, err := repository.GetService(db, userId, serviceId)
		if err != nil {
			util.Error(&w, http.StatusNoContent, err.Error())
			return
		}

		util.Send(&w, services)
	})
}

func createService(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user ID
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		// get service details from request
		var service dto.ServiceDto
		util.Receive[dto.ServiceDto](r.Body, &service)

		// call repository
		serviceId, err := repository.CreateService(db, userId, service)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
		}

		// return service ID
		util.Send(&w, serviceId)
	})
}

func cancelService(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))
		// get service param from query
		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			// if no service ID is found in query, return all services
            util.Error(&w, http.StatusBadRequest, "no service id")
            return
		}
		serviceId, err := strconv.ParseUint(serviceIdQuery, 10, 64)
		if err != nil {
			util.Error(&w, http.StatusBadRequest, err.Error())
			return
		}
		if err := repository.CancelService(db, userId, serviceId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
