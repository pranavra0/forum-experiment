package controllers

import (
	"net/http"
	"strconv"
	"forum-experiment/models"
	"github.com/go-chi/chi/v5"

)

func Home(w http.ResponseWriter, r *http.Request) {
	threads, err := models.GetAllThreads()
	if err != nil {
		http.Error(w, "unable to load threads", http.StatusInternalServerError)
		return
	}

	Render(w, "home", PageData{
		Name:    "home",
		Threads: threads,
	})
}

func NewThreadForm(w http.ResponseWriter, r *http.Request) {
	Render(w, "new", PageData{Name: "new"})
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

	_, err := models.CreateThread(title, content)
	if err != nil {
		http.Error(w, "could not create thread", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ShowThread(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	thread, err := models.GetThreadByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	replies, err := models.GetRepliesByThreadID(id)
	if err != nil {
		http.Error(w, "could not load replies", http.StatusInternalServerError)
		return
	}

	Render(w, "view_thread", PageData{
		Name:    "view_thread",
		Threads: []models.Thread{thread},
		Replies: replies,
	})
}

