package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

func main() {
	// create database connection
	db, err := sql.Open("pgx", "postgres://back:root@localhost:5432/back_test")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// create routes
	http.Handle("/user", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getUserHandler(db)))

	http.Handle("/user/update", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPut),
		middleware.Auth)(updateUserHandler(db)))

	http.Handle("/user/delete", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(deleteUserHandler(db)))

	if err = http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
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
		json.NewEncoder(w).Encode(user)
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
