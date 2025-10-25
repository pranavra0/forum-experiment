package controllers

import (
    "html/template"
    "log"
    "net/http"
)

var tmpl *template.Template

func init() {
    // Remove the err declaration since template.Must handles errors
    tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func Render(w http.ResponseWriter, name string, data PageData) {
    err := tmpl.ExecuteTemplate(w, "base", data)
    if err != nil {
        log.Printf("Template error: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}