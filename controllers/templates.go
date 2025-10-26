package controllers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
    "strings"
)

var templates map[string]*template.Template

func InitTemplates() {

	templates = make(map[string]*template.Template)

	funcs := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}

	files, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to glob templates: %v", err)
	}

    for _, file := range files {
        base := filepath.Base(file)
        name := strings.TrimSuffix(base, ".html")

		tmpl := template.New(name).Funcs(funcs)

        tmpl, err = tmpl.ParseFiles("templates/base.html", file)
        if err != nil {
            log.Fatalf("Failed to parse template %s: %v", name, err)
        }

        // Instead of storing all templates globally, store each page separately
        templates[name] = tmpl
    }

}

func Render(w http.ResponseWriter, name string, data map[string]any) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

