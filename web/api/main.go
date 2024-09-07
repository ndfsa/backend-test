package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/cardboard-bank/common/model"
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

	mdf := middleware.NewMiddlewareFactory(logger)

	srvhf := NewServicesHandlerFactory(repository.NewSrvRepository(db), mdf)
	usrhf := NewUsersHandlerFactory(repository.NewUsrRepository(db), mdf)

	jobQueue := make(chan model.Transaction, 2)
	tRepo := repository.NewTrsRepository(db, jobQueue)
	NewWorkerPool(5, jobQueue, logger)
	trshf := NewTransactionsHandlerFactory(tRepo, mdf)

	http.Handle("GET /users/{id}", usrhf.ReadSingleUser())
	http.Handle("GET /users", usrhf.ReadMultipleUsers())
	http.Handle("POST /users", usrhf.CreateUser())
	http.Handle("PUT /users", usrhf.UpdateUser())

    http.Handle("GET /users/{id}/services", srvhf.ReadUserServices())

	http.Handle("GET /services/{id}", srvhf.ReadSingleService())
	http.Handle("GET /services", srvhf.ReadMultipleServices())
	http.Handle("POST /services", srvhf.CreateService())
	http.Handle("PUT /services", srvhf.UpdateService())
	http.Handle("DELETE /services/{id}", srvhf.DeleteService())

    http.Handle("GET /services/{id}/transactions", trshf.ReadMultipleTransactions())

	http.Handle("GET /transactions/{id}", trshf.ReadSingleTransaction())
	http.Handle("GET /transactions", trshf.ReadMultipleTransactions())
	http.Handle("POST /transactions", trshf.CreateTransaction())

	logger.Println("---Starting API---")
	if err := http.ListenAndServe(":"+os.Getenv("AUTH_PORT"), nil); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
	close(jobQueue)
}
