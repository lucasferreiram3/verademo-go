package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func OpenDB() (*sql.DB, error) {
	var err error

	DB, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		return nil, err
	}

	return DB, nil
}
