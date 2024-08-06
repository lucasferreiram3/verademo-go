// Will possibly be the main driver file for the app
<<<<<<< HEAD
// App
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/view/", viewHandler)
	router.HandleFunc("/edit/", editHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
=======
package main

import (
	"fmt"
	"net/http"
)

func helloPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/", helloPage)
	http.ListenAndServe(" ", nil)
>>>>>>> 9d2378eb8e715ce6337363f977f5bcff34e3c573
}
