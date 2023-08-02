package repository

import "database/sql"

func AuthenticateUser(db *sql.DB, username string, password string) (uint64, error) {
	row := db.QueryRow("SELECT AUTHENTICATE_USER($1, $2)", username, password)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
