package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		middleware.Methods("POST", "GET"),
		middleware.Auth)(userHandler(db)))

	err = http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}

func userHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]

			userId, err := token.GetUserId(tokenString)
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
		case "POST":
			userId, err := repository.CreateUser(db, r.Body)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Could not create user", http.StatusInternalServerError)
				return
			}

			fmt.Fprint(w, userId)
		case "PUT":
            // TODO
			err := repository.UpdateUser()
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Could not update user", http.StatusInternalServerError)
				return
			}
		case "DELETE":
            // TODO
		}
	})
}
