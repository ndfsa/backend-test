package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/ndfsa/cardboard-bank/web/token"
)

const (
	CHECK_USER = iota
	CHECK_SERVICE
	CHECK_TRANSACTION
)

type MiddlewareFactory struct {
	db  *sql.DB
	log *log.Logger
}

func NewMiddlewareFactory(db *sql.DB, logger *log.Logger) MiddlewareFactory {
	return MiddlewareFactory{db, logger}
}

type (
	ErrHandler        = func(http.ResponseWriter, *http.Request) error
	RecoverMiddleware = func(ErrHandler) http.Handler
	ErrMiddleware     = func(ErrHandler) ErrHandler
)

func (factory *MiddlewareFactory) Auth(next ErrHandler) ErrHandler {
	return ErrHandler(func(w http.ResponseWriter, r *http.Request) error {
		encodedToken := r.Header.Get("Authorization")
		userId, err := token.ValidateAccessToken(encodedToken, token.KEY)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		return next(w, r.WithContext(ctx))
	})
}

// func (factory *MiddlewareFactory) Permissions(next http.Handler) http.Handler {
// 	if factory.db == nil {
// 		panic("db access is required for permissions middleware")
// 	}
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// }

func (factory *MiddlewareFactory) UploadLimit(limit int64) ErrMiddleware {
	return func(next ErrHandler) ErrHandler {
		return ErrHandler(func(w http.ResponseWriter, r *http.Request) error {
			if r.ContentLength > int64(limit) {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return errors.New("content length limit exceeded")
			}
			return next(w, r)
		})
	}
}

func (provider *MiddlewareFactory) Logger(next ErrHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provider.log.Printf("%s %s from %s %d bytes\n",
			r.Method,
            r.RequestURI,
			r.RemoteAddr,
			r.ContentLength)
		err := next(w, r)
		if err != nil {
			provider.log.Println(err)
		}
	})
}

func RecoverChain(
	recoverMiddleware RecoverMiddleware,
	errMiddlewares ...ErrMiddleware,
) RecoverMiddleware {
	return RecoverMiddleware(func(endpoint ErrHandler) http.Handler {
		if len(errMiddlewares) == 0 {
			return recoverMiddleware(endpoint)
		}
		next := endpoint
		for _, mid := range slices.Backward(errMiddlewares) {
			next = mid(next)
		}
		return recoverMiddleware(next)
	})
}
