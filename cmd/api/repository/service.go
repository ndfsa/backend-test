package repository

import (
	"database/sql"

	"github.com/ndfsa/backend-test/cmd/api/dto"
	"github.com/ndfsa/backend-test/internal/model"
	"github.com/shopspring/decimal"
)

func GetServices(db *sql.DB, userId uint64) ([]model.Service, error) {
	rows, err := db.Query(`SELECT * FROM GET_USER_SERVICES($1)`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var service model.Service
		if err := rows.Scan(
			&service.Id,
			&service.Type,
			&service.State,
			&service.Currency,
			&service.InitBalance,
			&service.Balance); err != nil {
			return services, err
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return services, err
	}

	return services, nil
}

func GetService(db *sql.DB, userId uint64, serviceId uint64) (model.Service, error) {
	rows := db.QueryRow("SELECT * FROM GET_USER_SERVICES($1) WHERE id = $2", userId, serviceId)

	var service model.Service

	if err := rows.Err(); err != nil {
		return service, err
	}

	if err := rows.Scan(
		&service.Id,
		&service.Type,
		&service.State,
		&service.Currency,
		&service.InitBalance,
		&service.Balance); err != nil {
		return service, err
	}

	return service, nil
}

func DebitService(db *sql.DB, userId uint64, serviceId uint64, amount decimal.Decimal) error {
	return nil
}

func CreditService(db *sql.DB, userId uint64, serviceId uint64, amount decimal.Decimal) error {
	return nil
}

func CreateService(db *sql.DB, userId uint64, service dto.ServiceDto) (uint64, error) {
	rows := db.QueryRow("SELECT CREATE_SERVICE($1, $2, $3, $4)",
		userId,
		service.Type,
		service.Currency,
		service.InitBalance)

	var serviceId uint64
	if err := rows.Err(); err != nil {
		return serviceId, err
	}

	if err := rows.Scan(&serviceId); err != nil {
		return serviceId, err
	}

	return serviceId, nil
}
