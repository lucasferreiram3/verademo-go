package routes

import (
	"net/http"
	"verademo-go/src-app/controllers"

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

func feedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowFeed(w, r)
	} else if r.Method == "POST" {

	}
}

func toolsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowTools(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessTools(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowLogin(w, r)
	} else if r.Method == "POST" {
		print("POST Reached")
		controllers.ProcessLogin(w, r)
	}

}

/*
Handler function used by router for register page.
*/
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowRegister(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessRegister(w, r)
	}
}

/*
Creates a router to listen for requests
Creates a session store
*/
func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feed", feedHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/tools", toolsHandler)
	router.HandleFunc("/register", registerHandler)
	router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	return router
}
