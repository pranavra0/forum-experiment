package controllers

import (
	"html/template"
	"log"
	"net/http"

	"forum-experiment/model"
)

var templates = template.Must(template.ParseFiles(
	"templates/base.html",
	"templates/home.html",
	"templates/new.html",
))

type PageData struct {
	Name     string
	Threads  []model.Thread
}

func Home(w http.ResponseWriter, r *http.Request) {
	threads, err := model.GetAllThreads()
	if err != nil {
		http.Error(w, "unable to load threads", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Name:    "home",
		Threads: threads,
	}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("template exec home:", err)
	}
}

func NewThreadForm(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Name: "new",
	}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("template exec new:", err)
	}
}


func CreateThread(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	if title == "" || content == "" {
		http.Error(w, "title and content required", http.StatusBadRequest)
		return
	}
	_, err := model.CreateThread(title, content)
	if err != nil {
		http.Error(w, "could not create thread", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
