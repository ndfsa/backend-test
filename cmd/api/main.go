package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// create database connection
	db, err := sql.Open("pgx", "postgres://back:root@localhost:5432/back_test")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

    CreateUserRoutes(db)
    CreateServiceRoutes(db)

	if err = http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}
