package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/cardboard-bank/common/repository"
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	repo := repository.NewAuthRepository(db)

	authf := NewAuthHandlerFactory(repo)

	http.Handle("POST /auth", authf.Authenticate())

	log.Println("---Starting AUTH---")
	if err := http.ListenAndServe(
		":"+os.Getenv("AUTH_PORT"), nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
