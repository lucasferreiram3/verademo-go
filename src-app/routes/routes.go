package routes

import (
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

// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Print("view")
// 	title := r.URL.Path[len("/view/"):]
// 	p, err := LoadPage(title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}
// 	Render(w, "view", p)
// }

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/edit/"):]
// 	p, err := LoadPage(title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}
// 	Render(w, "edit", p)
// 	// t, _ := template.ParseFiles("edit.html")
// 	// t.Execute(w, p)
// }

func feedHandler(w http.ResponseWriter, r *http.Request) {
	filename := "feed.html"
	if r.Method == "GET" {
		body, _ := os.ReadFile(filename)
		// if err != nil {
		// 	return nil, err
		// }
		Render(w, filename, body)
	} else if r.Method == "POST" {

	}
}

func Routes() *mux.Router {
	router := mux.NewRouter()
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	router.HandleFunc("/feed", feedHandler)
	// router.HandleFunc("/view/", viewHandler)
	// router.HandleFunc("/edit/", editHandler)

	return router
}

// Set up templates
var templates = template.Must(template.ParseGlob("templates/*.html"))

// htmlData is a byte array read from our template files
func Render(w http.ResponseWriter, filename string, htmlData []byte) {
	err := templates.ExecuteTemplate(w, filename, htmlData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func LoadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }
