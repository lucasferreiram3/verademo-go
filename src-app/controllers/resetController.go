package controllers

import (
	"log"
	"net/http"
	"os"
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

	if confirm != "Confirm" {
		errMsg := "Check the checkbox to confirm that you want to reset the database."
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}

	os.Remove("db.sqlite3")

	schema, err := os.ReadFile("blab_schema.sql")
	if err != nil {
		errMsg := "Failed to read schema file: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}

	// Open the new database
	sqlite.OpenDB()

	// Execute the schema to create tables and other database structures
	_, err = sqlite.DB.Exec(string(schema))
	if err != nil {
		errMsg := "Error executing schema: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "reset.html", outputs)
		return
	}

	view.Render(w, "reset.html", outputs)
}
