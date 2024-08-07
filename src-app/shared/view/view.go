package view

import (
	"net/http"
	"text/template"
)

// Set up templates
var templates = template.Must(template.ParseGlob("templates/*.html"))

// data is a struct with the variables to be passed into the template
func Render(w http.ResponseWriter, filename string, data []byte) {
	err := templates.ExecuteTemplate(w, filename, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
