package repository

import "database/sql"

func authenticateUser(db *sql.DB, username string, password string) (uint64, error) {
	row, err := db.Query("SELECT AUTHENTICATE_USER($1, $2)", username, password)
	if err != nil {
		return 0, err
	}
    var id uint64
    row.Scan(&id)

    return id, nil
}
