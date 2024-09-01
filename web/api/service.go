package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type ServicesHandlerFactory struct {
	repo repository.ServicesRepository
	mdf  middleware.MiddlewareFactory
}

func NewServicesHandlerFactory(
	repo repository.ServicesRepository,
	mdf middleware.MiddlewareFactory,
) ServicesHandlerFactory {
	return ServicesHandlerFactory{repo, mdf}
}

func (factory *ServicesHandlerFactory) CreateService() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateServiceRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		service, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}

		if err := factory.repo.CreateService(
			r.Context(), service); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(CreateServiceResponseDTO{
			service.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *ServicesHandlerFactory) ReadSingleService() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		serviceId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		service, err := factory.repo.GetService(r.Context(), serviceId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(service); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *ServicesHandlerFactory) ReadMultipleServices() http.Handler {
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

		services, err := factory.repo.GetServices(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := json.NewEncoder(w).Encode(services); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}

func (factory *ServicesHandlerFactory) CancelService() http.Handler {
	mid := middleware.RecoverChain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) error {
		serviceIdString := r.PathValue("id")
		if serviceIdString == "" {
			w.WriteHeader(http.StatusBadRequest)
			return errors.New("no id provided")
		}

		serviceId, err := uuid.Parse(serviceIdString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := factory.repo.SetServiceState(r.Context(), serviceId, "CLD"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
	return mid(f)
}
