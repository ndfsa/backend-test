package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/backend-test/cmd/back/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
)

func main() {
	db, err := sql.Open("pgx", "postgres://back:root@localhost:5432/back_test")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	http.Handle("/user", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getUserHandler(db)))

	http.Handle("/user/create", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPost))(createUserHandler(db)))

	http.Handle("/user/update", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPut),
		middleware.Auth)(updateUserHandler(db)))

	http.Handle("/user/delete", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(deleteUserHandler(db)))

	err = http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}

func getUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := token.GetUserId(r.Header.Get("Authorization"))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := repository.ReadUser(db, userId)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Could not read user", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, user)
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
		if err := repository.UpdateUser(db, r.Body); err != nil {
			log.Println(err.Error())
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	})
}

func deleteUserHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := repository.DeleteUser(); err != nil {
			log.Println(err.Error())
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
	})
}
