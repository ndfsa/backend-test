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

	repo := repository.NewAuthRepository(db)
	mdf := middleware.NewMiddlewareFactory(logger)

	authf := NewAuthHandlerFactory(repo, mdf)

	http.Handle("POST /auth", authf.Authenticate())

	logger.Println("---Starting AUTH---")
	if err := http.ListenAndServe(
		":"+os.Getenv("AUTH_PORT"), nil); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
}
