package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
)

func main() {
	db, err := sql.Open("pgx", "")
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

type UserDto struct {
	Name     string `json:"name"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

func userHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userId, err := token.GetUserId(strings.Split(r.Header.Get("Authorization"), " ")[1])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			fmt.Fprint(w, userId)
		case "POST":
			decoder := json.NewDecoder(r.Body)
			var newUser UserDto

			err := decoder.Decode(&newUser)

			if err != nil ||
				newUser.Name == "" ||
				newUser.Password == "" ||
				newUser.Username == "" {

				http.Error(w, "Malformed request", http.StatusBadRequest)
				return
			}
			var res uint64
			row := db.QueryRow("select create_user($1, $2, $3)",
				newUser.Name,
				newUser.Username,
				newUser.Password)
            err = row.Err()
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
                return
			}
			row.Scan(&res)
			fmt.Fprint(w, res)
		}
	})
}
