package repository

import "database/sql"

type SelfRepository struct {
	db *sql.DB
}

func NewSelfRepository(db *sql.DB) SelfRepository {
	return SelfRepository{db}
}


