package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	db "verademo/go/src-app/shared/db"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("view")
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func feedHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var database *sql.DB
	database, _ = db.InitDB()
	fmt.Print("s")
	router := mux.NewRouter()
	router.HandleFunc("/feed", feedHandler)
	router.HandleFunc("/view/", viewHandler)
	router.HandleFunc("/edit/", editHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
	database.Close()
}
