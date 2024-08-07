package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
}

/**
The routes file will create a Mux router, and then define all the
handlefunctions which will be run when a specified URL is received.
- Then, the handler functions can pass funcitonality to controllers,
so the controllers should have similar structure and be able to process a login, etc.

*/

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

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feed", feedHandler)
	router.HandleFunc("/view/", viewHandler)
	router.HandleFunc("/edit/", editHandler)

	return router
}
