package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	sqlite "verademo-go/src-app/shared/db"
	view "verademo-go/src-app/shared/view"
)

func ShowReset(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "reset.html", nil)
}

func ProcessReset(w http.ResponseWriter, r *http.Request) {
	type Error struct {
		Error string
	}

	confirm := r.FormValue("confirm")

	var outputs Error

	// Make sure user confirmed reset
	if confirm != "Confirm" {
		errMsg := "Check the checkbox to confirm that you want to reset the database."
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}

	// Get path of db folder
	// Get the current file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		errMsg := "Error getting current file."
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}

	// Get the path to the database folder
	dir := filepath.Join(filepath.Dir(currentFile), "..", "..", "db")

	// Delete the old database
	os.Remove(filepath.Join(dir, "db.sqlite3"))

	// Read the schema file
	schema, err := os.ReadFile(filepath.Join(dir, "blab_schema.sql"))
	if err != nil {
		errMsg := "Failed to read schema file: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}

	// Open the new database
	sqlite.OpenDB()

	// Execute the schema
	_, err = sqlite.DB.Exec(string(schema))
	if err != nil {
		errMsg := "Error executing schema: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}
	log.Println("Database reset successful.")
	view.Render(w, "reset.html", outputs)
}
