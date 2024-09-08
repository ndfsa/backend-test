package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type UsersHandlerFactory struct {
	repo repository.UsersRepository
}

func NewUsersHandlerFactory(
	repo repository.UsersRepository,
) UsersHandlerFactory {
	return UsersHandlerFactory{repo}
}

func (factory *UsersHandlerFactory) CreateUser() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000))
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateUserRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		user, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.CreateUser(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.CreateUserResponseDTO{
			Id: user.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *UsersHandlerFactory) ReadSingleUser() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		userIdString := r.PathValue("id")
		if userIdString == "" {
			w.WriteHeader(http.StatusNotFound)
			log.Println(errors.New("no id provided"))
			return
		}
		userId, err := uuid.Parse(userIdString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		user, err := factory.repo.FindUser(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.ReadUserResponseDTO{
			Id:       user.Id.String(),
			Role:     user.Role,
			Username: user.Username,
			Fullname: user.Fullname,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *UsersHandlerFactory) ReadMultipleUsers() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				return
			}
		} else {
			cursor = uuid.UUID{}
		}

		users, err := factory.repo.FindAllUsers(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res := make([]dto.ReadUserResponseDTO, 0, len(users))
		for _, user := range users {
			res = append(res, dto.ReadUserResponseDTO{
				Id:       user.Id.String(),
				Role:     user.Role,
				Username: user.Username,
				Fullname: user.Fullname,
			})
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *UsersHandlerFactory) UpdateUser() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.UpdateUserRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		user, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.UpdateUser(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}
