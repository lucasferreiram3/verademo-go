package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	var err error
	var db *sql.DB

	os.Create("db.sqlite3")

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
