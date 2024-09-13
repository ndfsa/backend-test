package main

import (
	"encoding/json"
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
	mdf  middleware.MiddlewareFactory
}

func NewServicesHandlerFactory(
	repo repository.ServicesRepository,
	mdf middleware.MiddlewareFactory,
) ServicesHandlerFactory {
	return ServicesHandlerFactory{repo, mdf}
}

func (factory *ServicesHandlerFactory) CreateService() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.Clearance(model.UserClearanceTeller))
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateServiceRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		service, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		ctx := r.Context()

		if err := factory.repo.CreateService(
			ctx, service); err != nil {
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

func (factory *ServicesHandlerFactory) UpdateUserService() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipUsr))
	f := func(w http.ResponseWriter, r *http.Request) {
        userId, _ := uuid.Parse(r.PathValue("id"))
        var req dto.UpdateUserServiceDTO
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

        if err := factory.repo.LinkServiceToUser(r.Context(), serviceId, userId); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Println(err)
            return
        }
    }
	return mid(http.HandlerFunc(f))
}

func (factory *ServicesHandlerFactory) ReadSingleService() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipSrv))
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceId, _ := uuid.Parse(r.PathValue("id"))
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
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.Clearance(model.UserClearanceTeller))
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

func (factory *ServicesHandlerFactory) UpdateService() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipSrv))
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceId, _ := uuid.Parse(r.PathValue("id"))
		var req dto.UpdateServiceRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipSrv))
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceId, _ := uuid.Parse(r.PathValue("id"))
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

func (factory *ServicesHandlerFactory) ReadUserServices() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipUsr))
	f := func(w http.ResponseWriter, r *http.Request) {
		userId, _ := uuid.Parse(r.PathValue("id"))
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

func (factory *ServicesHandlerFactory) CreateUserService() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipUsr))
	f := func(w http.ResponseWriter, r *http.Request) {
		userId, _ := uuid.Parse(r.PathValue("id"))
		var req dto.CreateServiceRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		service, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.CreateService(r.Context(), service); err != nil {
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
