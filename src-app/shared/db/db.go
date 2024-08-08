package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	var err error

	db, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		return nil, err
	}

	sqlStatement := `CREATE TABLE IF NOT EXISTS blabs (
						blabid INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
						blabber TEXT NOT NULL,
						content TEXT,
						timestamp TEXT
					);`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return nil, err
	}

	return db, nil
}
