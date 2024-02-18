package main

import (
	"log"
	"net/http"

	"github.com/ndfsa/cardboard-bank/cmd/api/dto"
	"github.com/ndfsa/cardboard-bank/cmd/api/repository"
	"github.com/ndfsa/cardboard-bank/internal/encoding"
	"github.com/ndfsa/cardboard-bank/internal/token"
)

func getUser(repo repository.UsersRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		// get user from database
		user, err := repo.Get(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// write user to response
		encoding.Send(w, user)
	})
}

func updateUser(repo repository.UsersRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}
		var user dto.UserDto
		if err := encoding.Receive(r, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := repo.Update(r.Context(), user, userId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}

func deleteUser(repo repository.UsersRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := r.Header.Get("Authorization")

		userId, err := token.GetUserId(encodedToken, tokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		// delete user from database
		if err := repo.Delete(r.Context(), userId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}
