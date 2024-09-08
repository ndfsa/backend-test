package middleware

import (
	"context"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/web/token"
)

const (
	USER_CTX_KEY = "USER"
	ERR_CTX_KEY  = "USER"

	logReset   = "\033[0m"
	logRed     = "\033[31m"
	logGreen   = "\033[32m"
	logYellow  = "\033[33m"
	logBlue    = "\033[34m"
	logMagenta = "\033[35m"
	logCyan    = "\033[36m"
	logGray    = "\033[37m"
	logWhite   = "\033[97m"
)

type Middleware = func(http.Handler) http.Handler

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")
		userId, err := token.ValidateAccessToken(encodedToken, token.KEY)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		ctx := context.WithValue(r.Context(), USER_CTX_KEY, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Clearence(level int8) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(USER_CTX_KEY).(model.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				log.Println("user not provided")
				return
			}
			if user.Role < level {
				w.WriteHeader(http.StatusForbidden)
				log.Println("clearence level too low")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Ownership(condition func() bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("try ownership now")
			if true {
				log.Println("fail on purpose")
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
				log.Println("content length limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s from %s %d[bytes] %f[s]",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.ContentLength,
			time.Since(start).Seconds())
	})
}

func Chain(middlewares ...Middleware) Middleware {
	fn := func(endpoint http.Handler) http.Handler {
		if len(middlewares) == 0 {
			return endpoint
		}

		next := endpoint
		for _, m := range slices.Backward(middlewares) {
			next = m(next)
		}

		return next
	}
	return Middleware(fn)
}
