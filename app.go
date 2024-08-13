package main

import (
	"embed"
	"log"
	"net/http"
	router "verademo-go/src-app/routes"
	sqlite "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	"verademo-go/src-app/shared/view"
)

//go:embed resources
var resources embed.FS

//go:embed templates
var templates embed.FS

func main() {

	session.Configure(session.Session{Name: "verademo", SecretKey: "key"})
	sqlite.OpenDB()
	log.Print("\nStarting VerademoGO....")
	log.Print("\nVerademoGO is running.")
	view.ParseTemplates(templates)
	router.SetResources(resources)
	log.Fatal(http.ListenAndServe(":8000", router.Routes()))
}
