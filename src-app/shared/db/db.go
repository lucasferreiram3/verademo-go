package db

import (
	"database/sql"
	"os"
)

func InitDB() (*sql.DB, error) {
	var err error
	var db *sql.DB

	os.Create("verademo/go/src-app/shared/db/db.sqlite3")

	db, err = sql.Open("sqlite3", "verademo/go/src-app/shared/db")
	if err != nil {
		return nil, err
	}

	sqlStatement := `CREATE TABLE IF NOT EXISTS blabs (
						blabid INTEGER NOT NULL AUTOINCREMENT PRIMARY KEY,
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
