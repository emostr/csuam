package db

import (
	"database/sql"
)

type Queries struct {
	db *sql.DB
}

func New(sdb *sql.DB) *Queries {
	return &Queries{db: sdb}
}
