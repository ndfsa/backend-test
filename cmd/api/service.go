package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

func CreateServiceRoutes(db *sql.DB, baseUrl string) {
	http.Handle(baseUrl+"/services", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getServices(db)))
	http.Handle(baseUrl+"/service", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getService(db)))
}
func getServices(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user id from header
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		// get all user services
		services, err := repository.GetServices(db, userId)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		util.Send(&w, services)
	})
}

func getService(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user id from header
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		// get service param from query
		serviceIdQuery := r.URL.Query().Get("id")
		if serviceIdQuery == "" {
			util.Error(&w, http.StatusBadRequest, "no id present")
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
