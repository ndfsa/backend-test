package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type UsersHandlerFactory struct {
	repo repository.UsersRepository
	mdf  middleware.MiddlewareFactory
}

func NewUsersHandlerFactory(
	repo repository.UsersRepository,
	mdf middleware.MiddlewareFactory,
) UsersHandlerFactory {
	return UsersHandlerFactory{repo, mdf}
}

func (factory *UsersHandlerFactory) CreateUser() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000))
	f := func(w http.ResponseWriter, r *http.Request) error {
		var req dto.CreateUserRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		user, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		if err := factory.repo.CreateUser(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(dto.CreateUserResponseDTO{
			Id: user.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *UsersHandlerFactory) ReadSingleUser() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		userIdString := r.PathValue("id")
		if userIdString == "" {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("no id provided")
		}
		userId, err := uuid.Parse(userIdString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		user, err := factory.repo.FindUser(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(dto.ReadUserResponseDTO{
			Id:       user.Id.String(),
			Role:     user.Role,
			Username: user.Username,
			Fullname: user.Fullname,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *UsersHandlerFactory) ReadMultipleUsers() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return err
			}
		} else {
			cursor = uuid.UUID{}
		}

		users, err := factory.repo.FindAllUsers(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
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
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *UsersHandlerFactory) UpdateUser() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		var req dto.UpdateUserRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		user, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		if err := factory.repo.UpdateUser(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}
