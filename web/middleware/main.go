package middleware

import (
	"context"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/token"
)

const (
	userKey      = "USER"
	clearanceKey = "CLEARANCE"

	logReset   = "\033[0m"
	logRed     = "\033[31m"
	logGreen   = "\033[32m"
	logYellow  = "\033[33m"
	logBlue    = "\033[34m"
	logMagenta = "\033[35m"
	logCyan    = "\033[36m"
	logGray    = "\033[37m"
	logWhite   = "\033[97m"

	OwnershipUsr = 'U'
	OwnershipSrv = 'S'
	OwnershipTrs = 'T'
)

type MiddlewareFactory struct {
	repo repository.OwnershipRepository
}

func NewMiddlewareFactory(repo repository.OwnershipRepository) MiddlewareFactory {
	return MiddlewareFactory{repo}
}

type Middleware = func(http.Handler) http.Handler

func (factory *MiddlewareFactory) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")
		userId, err := token.ValidateAccessToken(encodedToken, token.KEY)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuthenticatedUser(ctx context.Context) model.User {
	return ctx.Value(userKey).(model.User)
}

func (factory *MiddlewareFactory) Clearance(level int8) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getAuthenticatedUser(r.Context())
			if user.Clearance < level {
				w.WriteHeader(http.StatusForbidden)
				log.Println("user does not have sufficient clearance")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (factory *MiddlewareFactory) ClearanceOrOwnership(level int8, entity rune) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getAuthenticatedUser(r.Context())
			if user.Clearance < level {
				resource, err := uuid.Parse(r.PathValue("id"))
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					log.Println(err)
					return
				}

				ctx := r.Context()

				var cerr error
				switch entity {
				case OwnershipUsr:
					cerr = factory.repo.CheckUserOwnership(resource, user.Id)
				case OwnershipSrv:
					cerr = factory.repo.CheckServiceOwnership(ctx, resource, user.Id)
				case OwnershipTrs:
					cerr = factory.repo.CheckTransactionOwnership(ctx, resource, user.Id)
				default:
					panic("unknown ownership entity")
				}

				if cerr != nil {
					w.WriteHeader(http.StatusForbidden)
					log.Println(cerr)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (factory *MiddlewareFactory) UploadLimit(limit int64) Middleware {
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

func (factory *MiddlewareFactory) Logger(next http.Handler) http.Handler {
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
