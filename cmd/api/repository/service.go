package repository

import (
	"database/sql"

	"github.com/ndfsa/backend-test/internal/model"
)

func GetServices(db *sql.DB) ([]model.Service, error) {
    db.Query("SELECT * FROM services")
	return nil, nil
}
