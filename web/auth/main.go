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
	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	authRepo := repository.NewAuthRepository(db)
    ownRepo := repository.NewOwnershipRepository(db)

    mdf := middleware.NewMiddlewareFactory(ownRepo)
	authf := NewAuthHandlerFactory(authRepo, mdf)

	http.Handle("POST /auth", authf.Authenticate())
	http.Handle("GET /refresh", authf.RefreshToken())

	log.Println("---Starting AUTH---")
	if err := http.ListenAndServe(
		":"+os.Getenv("AUTH_PORT"), nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
