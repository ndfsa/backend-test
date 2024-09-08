package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type ServicesHandlerFactory struct {
	repo repository.ServicesRepository
}

func NewServicesHandlerFactory(
	repo repository.ServicesRepository,
) ServicesHandlerFactory {
	return ServicesHandlerFactory{repo}
}

func (factory *ServicesHandlerFactory) CreateService() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateServiceRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		service, userId, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.CreateService(
			r.Context(), service); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := factory.repo.LinkServiceToUser(r.Context(), service.Id, userId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.CreateServiceResponseDTO{
			Id: service.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *ServicesHandlerFactory) ReadSingleService() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		service, err := factory.repo.FindService(r.Context(), serviceId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.ReadServiceResponseDTO{
			Id:          service.Id.String(),
			Type:        service.Type,
			State:       service.State,
			Currency:    service.Currency,
			InitBalance: service.InitBalance.String(),
			Balance:     service.Balance.String(),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *ServicesHandlerFactory) ReadMultipleServices() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth,
		middleware.Clearence(model.UserRoleTeller))
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

		services, err := factory.repo.FindAllServices(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res := make([]dto.ReadServiceResponseDTO, 0, len(services))
		for _, service := range services {
			res = append(res, dto.ReadServiceResponseDTO{
				Id:          service.Id.String(),
				Type:        service.Type,
				State:       service.State,
				Currency:    service.Currency,
				InitBalance: service.InitBalance.String(),
				Balance:     service.Balance.String(),
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

func (factory *ServicesHandlerFactory) ReadUserServices() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		userId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}

		services, err := factory.repo.FindUserServices(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res := make([]dto.ReadServiceResponseDTO, 0, len(services))
		for _, service := range services {
			res = append(res, dto.ReadServiceResponseDTO{
				Id:          service.Id.String(),
				Type:        service.Type,
				State:       service.State,
				Currency:    service.Currency,
				InitBalance: service.InitBalance.String(),
				Balance:     service.Balance.String(),
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

func (factory *ServicesHandlerFactory) UpdateService() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)

	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.UpdateServiceRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		serviceId, err := uuid.Parse(req.Id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.UpdateService(r.Context(), model.Service{
			Id:    serviceId,
			State: req.State,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	return mid(http.HandlerFunc(f))
}

func (factory *ServicesHandlerFactory) DeleteService() http.Handler {
	mid := middleware.Chain(
		middleware.Logger,
		middleware.UploadLimit(1000),
		middleware.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceIdString := r.PathValue("id")
		if serviceIdString == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(errors.New("no id provided"))
			return
		}

		serviceId, err := uuid.Parse(serviceIdString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := factory.repo.UpdateService(r.Context(), model.Service{
			Id:    serviceId,
			State: model.ServiceStateClosed,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}
