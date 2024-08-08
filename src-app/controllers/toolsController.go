package controllers

import (
	"log"
	"net/http"
	"verademo-go/src-app/shared/view"
)

func ShowTools(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "tools.html", nil)
}

func ProcessTools(w http.ResponseWriter, r *http.Request) {
	host := r.FormValue("host")
	fortuneFile := r.FormValue("fortunefile")

	if host == "" {
		ping := ""
	} else {
		ping := ping(host)
	}

	log.Print(host)
	view.Render(w, "tools.html")
}

func ping(host string) string {
	return host
}
