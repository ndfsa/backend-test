package main

import (
	"database/sql"
	"net/http"

	"github.com/ndfsa/backend-test/cmd/api/repository"
	"github.com/ndfsa/backend-test/internal/middleware"
	"github.com/ndfsa/backend-test/internal/token"
	"github.com/ndfsa/backend-test/internal/util"
)

func CreateUserRoutes(db *sql.DB, baseUrl string) {
	repo := repository.NewUsersRepository(db)

	http.Handle(baseUrl+"/user", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodGet),
		middleware.Auth)(getUserHandler(repo)))

	http.Handle(baseUrl+"/user/update", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodPut),
		middleware.Auth)(updateUserHandler(repo)))

	http.Handle(baseUrl+"/user/delete", middleware.Chain(
		middleware.Logger,
		middleware.Method(http.MethodDelete),
		middleware.Auth)(deleteUserHandler(repo)))
}

func getUserHandler(repo repository.UsersRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r)

		// get user from database
		user, err := repo.Get(r.Context(), userId)
		if err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}

		// write user to response
		util.Send(&w, user)
	})
}

func updateUserHandler(repo repository.UsersRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r)

		if err := repo.Update(r.Context(), r.Body, userId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

func deleteUserHandler(repo repository.UsersRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open jwt to retrieve userId
		userId, _ := token.GetUserId(r)

		// delete user from database
		if err := repo.Delete(r.Context(), userId); err != nil {
			util.Error(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})
}
