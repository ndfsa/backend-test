package middleware

import (
	"log"
	"net/http"

	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

type Middleware = func(http.Handler) http.Handler

var Basic = Chain(Logger, UploadLimit(1000))

func Auth(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encodedToken := r.Header.Get("Authorization")
			if err := token.ValidateAccessToken(encodedToken, key); err != nil {
				util.SendError(&w, http.StatusUnauthorized, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UploadLimit(limit int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > int64(limit) {
				util.SendError(&w, http.StatusRequestEntityTooLarge, "request too large")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] from %s of size %d\n", r.Method, r.RemoteAddr, r.ContentLength)
		next.ServeHTTP(w, r)
	})
}

func Chain(middlewares ...Middleware) Middleware {
	fn := func(endpoint http.Handler) http.Handler {
		if len(middlewares) == 0 {
			return endpoint
		}
		next := middlewares[len(middlewares)-1](endpoint)
		for i := len(middlewares) - 2; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
	return Middleware(fn)
}
