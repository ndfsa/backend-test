package middleware

import (
	"log"
	"net/http"

	"github.com/ndfsa/cardboard-bank/internal/token"
)

type Middleware = func(http.Handler) http.Handler

var Basic = Chain(Logger, UploadLimit(1000))

func BasicAuth(key string) Middleware {
	return Chain(Logger, Auth(key), UploadLimit(1000))
}

func Auth(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encodedToken := r.Header.Get("Authorization")
			if err := token.ValidateAccessToken(encodedToken, key); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Println(err)
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
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				log.Println("request too large")
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
