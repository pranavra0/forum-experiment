package controllers

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

// helper to render a named page template with data
func Render(w http.ResponseWriter, name string, data any) {
	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("template exec %s: %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
