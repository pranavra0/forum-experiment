package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

var templates map[string]*template.Template

func InitTemplates() {

	templates = make(map[string]*template.Template)

	funcs := template.FuncMap{
		"add":         func(a, b int) int { return a + b },
		"sub":         func(a, b int) int { return a - b },
		"highlight":   highlight,
		"dict":        dict,
		"mul":         mul,
		"mod":         mod,
		"parseQuotes": parseQuotes,
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

func highlight(text, query string) template.HTML {
	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(query))
	highlighted := re.ReplaceAllStringFunc(text, func(match string) string {
		return `<mark style="background-color: #fff176; color: black; padding: 0 2px; border-radius: 3px;">` + match + `</mark>`
	})
	return template.HTML(highlighted)
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call: uneven number of arguments")
	}
	d := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings, got %T", values[i])
		}
		d[key] = values[i+1]
	}
	return d, nil
}

func mul(a, b int) int {
	return a * b
}

func mod(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

func parseQuotes(content string) template.HTML {
	re := regexp.MustCompile(`(?s)\[quote=(.*?)\](.*?)\[/quote\]`)

	// While there is at least one [quote], replace it
	for re.MatchString(content) {
		content = re.ReplaceAllStringFunc(content, func(match string) string {
			parts := re.FindStringSubmatch(match)
			if len(parts) < 3 {
				return match
			}
			user := template.HTMLEscapeString(parts[1])
			inner := parseQuotes(parts[2]) // <-- recurse here!
			return `<div class="quote"><strong>` + user + ` said:</strong>` + string(inner) + `</div>`
		})
	}

	return template.HTML(content)
}
