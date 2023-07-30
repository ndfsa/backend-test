package middleware

import (
	"log"
	"net/http"

	"golang.org/x/exp/slices"

	"github.com/ndfsa/backend-test/auth/token"
)

type middleware = func(http.Handler) http.Handler

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			err := token.Validate(w, r)

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

func Methods(supportedMethods ...string) middleware {
	return middleware(

		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				if !slices.Contains(supportedMethods, r.Method) {
					http.Error(w, "method unsupported", http.StatusMethodNotAllowed)
					log.Printf("middleware - Methods: %s not in supported methods %v\n",
						r.Method,
						supportedMethods)
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
