package view

import (
	"net/http"
	"text/template"
)

// Set up templates
var templates = template.Must(template.ParseGlob("templates/*.html"))

// htmlData is a byte array read from our template files
func Render(w http.ResponseWriter, filename string, htmlData []byte) {
	err := templates.ExecuteTemplate(w, filename, htmlData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
