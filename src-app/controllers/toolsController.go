package controllers

import (
	"net/http"
	"verademo-go/src-app/shared/view"
)

func ShowTools(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "tools.html", nil)
}
