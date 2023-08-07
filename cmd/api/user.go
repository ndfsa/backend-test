package main

import (
	"database/sql"
	"net/http"

	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

func CreateUserRoutes(db *sql.DB) {
	http.Handle("/api/user", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getUserHandler(db)))

	http.Handle("/api/user/update", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPut),
		middleware.Auth)(updateUserHandler(db)))

	http.Handle("/api/user/delete", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(deleteUserHandler(db)))
}

func getUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		// get user from database
		user, err := repository.ReadUser(db, userId)
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
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		if err := repository.UpdateUser(db, r.Body, userId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

func deleteUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r.Header.Get("Authorization"))

		// delete user from database
		if err := repository.DeleteUser(db, userId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}
