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
	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	jobQueue := make(chan model.Transaction, 2)

	usrRepo := repository.NewUsrRepository(db)
	srvRepo := repository.NewSrvRepository(db)
	trsRepo := repository.NewTrsRepository(db, jobQueue)
	ownRepo := repository.NewOwnershipRepository(db)

	mdf := middleware.NewMiddlewareFactory(ownRepo)

	srvhf := NewServicesHandlerFactory(srvRepo, mdf)
	usrhf := NewUsersHandlerFactory(usrRepo, mdf)
	trshf := NewTransactionsHandlerFactory(trsRepo, mdf, srvRepo)

	NewWorkerPool(5, jobQueue)

	http.Handle("GET /users/{id}", usrhf.ReadSingleUser())
	http.Handle("GET /users", usrhf.ReadMultipleUsers())
	http.Handle("POST /users", usrhf.CreateUser())
	http.Handle("PUT /users/{id}", usrhf.UpdateUser())

	http.Handle("GET /users/{id}/services", srvhf.ReadUserServices())
	http.Handle("POST /users/{id}/services", srvhf.CreateUserService())
	http.Handle("PUT /users/{id}/services", srvhf.UpdateUserService())

	http.Handle("GET /services/{id}", srvhf.ReadSingleService())
	http.Handle("GET /services", srvhf.ReadMultipleServices())
	http.Handle("POST /services", srvhf.CreateService())
	http.Handle("PUT /services/{id}", srvhf.UpdateService())
	http.Handle("DELETE /services/{id}", srvhf.DeleteService())

	http.Handle("GET /services/{id}/transactions", trshf.ReadServiceTransactions())

	http.Handle("GET /transactions/{id}", trshf.ReadSingleTransaction())
	http.Handle("GET /transactions", trshf.ReadMultipleTransactions())
	http.Handle("POST /transactions", trshf.CreateTransaction())

	log.Println("---Starting API---")
	if err := http.ListenAndServe(":"+os.Getenv("AUTH_PORT"), nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	close(jobQueue)
}
