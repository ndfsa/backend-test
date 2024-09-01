package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

func main() {
	logger := log.Default()

	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	mdf := middleware.NewMiddlewareFactory(db, logger)

	srvhf := NewServicesHandlerFactory(repository.NewServiceRepository(db), mdf)
	trshf := NewTransactionsHandlerFactory(repository.NewTransactionsRepository(db, 0, 2), mdf)
	usrhf := NewUsersHandlerFactory(repository.NewUserRepository(db), mdf)

	http.Handle("GET /services/{id}", srvhf.ReadSingleService())
	http.Handle("GET /services", srvhf.ReadMultipleServices())
	http.Handle("POST /services", srvhf.CreateService())
	http.Handle("DELETE /services/{id}", srvhf.CancelService())

	http.Handle("GET /users/{id}", usrhf.ReadSingleUser())
	http.Handle("GET /users", usrhf.ReadMultipleUsers())
	http.Handle("POST /users", usrhf.CreateUser())
	http.Handle("PUT /users", usrhf.UpdateUser())

	http.Handle("POST /transactions", trshf.CreateTransaction())
	http.Handle("GET /transactions/{id}", trshf.ReadSingleTransaction())
	http.Handle("GET /transactions", trshf.ReadMultipleTransactions())

	logger.Println("---Starting API---")
	if err := http.ListenAndServe(":"+os.Getenv("AUTH_PORT"), nil); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
}
