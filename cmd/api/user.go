package main

import (
	"database/sql"
	"net/http"

	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

func CreateUserRoutes(db *sql.DB, baseUrl string) {
	http.Handle(baseUrl+"/user", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getUserHandler(db)))

	http.Handle(baseUrl+"/user/update", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPut),
		middleware.Auth)(updateUserHandler(db)))

	http.Handle(baseUrl+"/user/delete", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(deleteUserHandler(db)))
}

func getUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r)

		// get user from database
		user, err := repository.ReadUser(r.Context(), db, userId)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		// write user to response
		util.Send(&w, user)
	})
}

func updateUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r)

		if err := repository.UpdateUser(r.Context(), db, r.Body, userId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

func deleteUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r)

		// delete user from database
		if err := repository.DeleteUser(r.Context(), db, userId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}
