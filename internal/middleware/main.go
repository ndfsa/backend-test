package middleware

import (
	"log"
	"net/http"

	"github.com/ndfsa/backend-test/internal/token"
)

type middleware = func(http.Handler) http.Handler

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			err := token.Validate(header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				log.Printf("middleware - Auth: %s\n", err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
}

func UploadLimit(limit int64) middleware {
	return middleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ContentLength > int64(limit) {
					http.Error(w, "Request body is too long", http.StatusRequestEntityTooLarge)
					log.Printf("middleware - UploadLimit: body is too long\n")
					return
				}
				next.ServeHTTP(w, r)
			})
		})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[%s] from %s of size %d\n", r.Method, r.RemoteAddr, r.ContentLength)
			next.ServeHTTP(w, r)
		})
}

func Method(method string) middleware {
	return middleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if method != r.Method {
					http.Error(w, "method unsupported", http.StatusMethodNotAllowed)

					log.Printf("middleware - Method: %s not in supported\n",
						r.Method)
					return
				}
				next.ServeHTTP(w, r)
			})
		})
}

func Chain(middlewares ...middleware) middleware {
	return middleware(
		func(endpoint http.Handler) http.Handler {
			if len(middlewares) == 0 {
				return endpoint
			}
			next := middlewares[len(middlewares)-1](endpoint)
			for i := len(middlewares) - 2; i >= 0; i-- {
				next = middlewares[i](next)
			}
			return next
		})
}
