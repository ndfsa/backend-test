package repository

import (
	"database/sql"

	"github.com/ndfsa/backend-test/cmd/api/dto"
	"github.com/ndfsa/backend-test/internal/model"
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

func CreateService(db *sql.DB, userId uint64, service dto.ServiceDto) (uint64, error) {
	idRow := db.QueryRow(`INSERT INTO services (type, state, currency, init_balance, balance)
        VALUES (&1, 'REQ', &2, $3, 0)
        RETURNING id`,
		service.Type,
		service.Currency,
		service.InitBalance)

	var serviceId uint64
	if err := idRow.Err(); err != nil {
		return 0, err
	}

	if err := idRow.Scan(&serviceId); err != nil {
		return 0, err
	}

	rows := db.QueryRow("INSERT INTO user_service (user_id, service_id) VALUES ($1, $2)",
		userId,
		serviceId)

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return serviceId, nil
}

func CancelService(db *sql.DB, userId uint64, serviceId uint64) error {
	rows := db.QueryRow(`UPDATE services SET state = 'CLD'
        FROM users JOIN user_service ON users.id = user_id
        WHERE users.id = $1 AND services.id = $2`,
		userId,
		serviceId)

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
