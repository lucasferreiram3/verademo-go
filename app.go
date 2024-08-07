package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	router "verademo/go/src-app/routes"
	db "verademo/go/src-app/shared/db"
)

func main() {
	var database *sql.DB
	database, _ = db.InitDB()
	fmt.Print("s")

	log.Fatal(http.ListenAndServe(":8080", router.Routes()))
	database.Close()
}
