package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/backend-test/cmd/api/repository"
	ilib "github.com/ndfsa/backend-test/internal"
)

func main() {
	// create database connection
	db, err := sql.Open("pgx", "postgres://back:root@localhost:5432/back_test")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// create routes
	http.Handle("/user", ilib.Chain(
		ilib.Logger,
		ilib.Method(http.MethodGet),
		ilib.Auth)(getUserHandler(db)))

	http.Handle("/user/create", ilib.Chain(
		ilib.Logger,
		ilib.Method(http.MethodPost))(createUserHandler(db)))

	http.Handle("/user/update", ilib.Chain(
		ilib.Logger,
		ilib.Method(http.MethodPut),
		ilib.Auth)(updateUserHandler(db)))

	http.Handle("/user/delete", ilib.Chain(
		ilib.Logger,
		ilib.Method(http.MethodDelete),
		ilib.Auth)(deleteUserHandler(db)))

	if err = http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}

func getUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // open jwt to retrieve userId
		userId, _ := ilib.GetUserId(r.Header.Get("Authorization"))

        // get user from database
		user, err := repository.ReadUser(db, userId)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Could not read user", http.StatusInternalServerError)
			return
		}

        // write user to response
		json.NewEncoder(w).Encode(user)
	})
}

func createUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := repository.CreateUser(db, r.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, userId)
	})
}

func updateUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // open jwt to retrieve userId
		userId, _ := ilib.GetUserId(r.Header.Get("Authorization"))

		if err := repository.UpdateUser(db, r.Body, userId); err != nil {
			log.Println(err.Error())
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	})
}

func deleteUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // open jwt to retrieve userId
		userId, _ := ilib.GetUserId(r.Header.Get("Authorization"))

        // delete user from database
		if err := repository.DeleteUser(db, userId); err != nil {
			log.Println(err.Error())
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
	})
}
