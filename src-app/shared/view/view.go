package view

import (
	"embed"
	"net/http"
	"text/template"
)

// Set up templates
var templates *template.Template

func ParseTemplates(t embed.FS) {
	templates = template.Must(template.ParseFS(t, "templates/*.html"))
}

// data is a struct with the variables to be passed into the template
func Render(w http.ResponseWriter, filename string, data any) {
	err := templates.ExecuteTemplate(w, filename, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
