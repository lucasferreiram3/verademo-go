package routes

import (
	"embed"
	"net/http"
	"path/filepath"
	"runtime"
	"verademo-go/src-app/controllers"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
}

// Set the embedded files
var resources embed.FS

func SetResources(r embed.FS) {
	resources = r
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
		controllers.ProcessFeed(w, r)
	}
}

func moreFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.MoreFeed(w, r)
	}
}

func blabHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowBlab(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessBlab(w, r)
	}
}

func blabbersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowBlabbers(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessBlabbers(w, r)
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

func passwordHintHandler(w http.ResponseWriter, r *http.Request) {
	controllers.ShowPasswordHint(w, r)
}
func totpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowTotp(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessTotp(w, r)
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

func registerFinishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowRegisterFinish(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessRegisterFinish(w, r)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowProfile(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessProfile(w, r)
	}
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controllers.ShowReset(w, r)
	} else if r.Method == "POST" {
		controllers.ProcessReset(w, r)
	}
}

/*
Creates a router to listen for requests
Creates a session store
*/
func Routes() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/", loginHandler)
	router.HandleFunc("/feed", feedHandler)
	router.HandleFunc("/morefeed", moreFeedHandler)
	router.HandleFunc("/blab", blabHandler)
	router.HandleFunc("/blabbers", blabbersHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/totp", totpHandler)
	router.HandleFunc("/logout", controllers.ProcessLogout)
	router.HandleFunc("/tools", toolsHandler)
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/reset", resetHandler)
	router.HandleFunc("/register-finish", registerFinishHandler)
	router.HandleFunc("/password-hint", passwordHintHandler)
	router.HandleFunc("/profile", profileHandler)
	router.HandleFunc("/downloadprofileimage", controllers.DownloadImage)
	router.PathPrefix("/resources/").Handler(http.FileServer(http.FS(resources)))
	_, currentFile, _, _ := runtime.Caller(0)
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(filepath.Dir(currentFile), "..", "..", "images")))))
	http.Handle("/", router)

	return router
}
