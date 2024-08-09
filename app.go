package main

import (
	"database/sql"
	"log"
	"net/http"
	router "verademo-go/src-app/routes"
	db "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
)

func main() {
	session.Configure(session.Session{Name: "verademo", SecretKey: "key"})
	var database *sql.DB
	database, _ = db.InitDB()
	log.Print("\nStarting VerademoGO....")
	log.Print("\nVerademoGO is running.")

	log.Fatal(http.ListenAndServe(":8080", router.Routes()))
	database.Close()
}
