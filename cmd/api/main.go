package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	tokenKey = "test-application"
	baseUrl  = "/api/v1"
)

func main() {
	// create database connection
	db, err := sql.Open("pgx", "postgres://back:root@localhost:5432/cardboard_bank")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	CreateUserRoutes(db, baseUrl, tokenKey)
	CreateServiceRoutes(db, baseUrl, tokenKey)

	if err = http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}
