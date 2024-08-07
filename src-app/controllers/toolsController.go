package controllers

import (
	"net/http"
	"os"
	"verademo-go/src-app/shared/view"
)

func ShowTools(w http.ResponseWriter, r *http.Request) {
	filename := "tools.html"
	body, _ := os.ReadFile(filename)
	view.Render(w, "tools.html", body)
}
