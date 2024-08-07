package controllers

import (
	"net/http"
	router "verademo/go/src-app/routes"
)

func showFeed(w http.ResponseWriter, r *http.Request) {
	p, err := router.LoadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	router.Render(w, "edit", p)
}
