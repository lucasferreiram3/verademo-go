package db

import (
	"database/sql"
	"log"
	"path/filepath"
	"runtime"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var DB *sql.DB

func OpenDB() (*sql.DB, error) {
	var err error

	// Get the current file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		errMsg := "Error getting current file."
		log.Println(errMsg)
		return nil, err
	}

	// Get the path to the database folder
	dir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "db")

	// Open the database
	DB, err = sql.Open("sqlite3", filepath.Join(dir, "db.sqlite3"))
	if err != nil {
		errMsg := "Error opening database: \n" + err.Error()
		log.Println(errMsg)
		return nil, err
	}

	return DB, nil
}
