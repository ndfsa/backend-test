package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/web/token"
)

const (
	USER_CTX_KEY = "USER"

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

type MiddlewareFactory struct {
	log *log.Logger
}

func NewMiddlewareFactory(logger *log.Logger) MiddlewareFactory {
	return MiddlewareFactory{logger}
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

		ctx := context.WithValue(r.Context(), USER_CTX_KEY, userId)
		return next(w, r.WithContext(ctx))
	})
}

func (factory *MiddlewareFactory) Clearence(level int8) ErrMiddleware {
	return func(next ErrHandler) ErrHandler {
		return ErrHandler(func(w http.ResponseWriter, r *http.Request) error {
			user, ok := r.Context().Value(USER_CTX_KEY).(model.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return errors.New("user not provided")
			}
			if user.Role < level {
				w.WriteHeader(http.StatusForbidden)
				return errors.New("clearence level too low")
			}
			return next(w, r)
		})
	}
}

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

func (factory *MiddlewareFactory) Logger(next ErrHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		err := next(w, r)
		info := fmt.Sprintf("%s %s from %s %d[bytes] %f[s]",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.ContentLength,
			time.Since(start).Seconds())
		if err != nil {
            factory.log.Printf("%s%s err: %s%s\n", logRed, info, err, logReset)
		} else {
			factory.log.Println(info)
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
