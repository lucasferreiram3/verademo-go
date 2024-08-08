package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	router "verademo-go/src-app/routes"
	db "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
)

type person struct{}

func main() {
	session.Configure(session.Session{Name: "verademo", SecretKey: "key"})
	var database *sql.DB
	database, _ = db.InitDB()

	fmt.Println("s")

	log.Fatal(http.ListenAndServe(":8080", router.Routes()))
	database.Close()
}
