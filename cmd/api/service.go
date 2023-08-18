package main

import (
	"database/sql"
	"net/http"

	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/util"
)

func CreateServiceRoutes(db *sql.DB) {
	http.HandleFunc("/api/service/", getService(db))
}
func getService(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        services, err := repository.GetServices(db)
        if err != nil {
            util.Error(&w, http.StatusInternalServerError, "could not get services from database")
        }

        util.Send(&w, services)
	})
}
