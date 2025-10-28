package controllers

import (
	"net/http"
	"strconv"

	"forum-experiment/models"

	"github.com/go-chi/chi/v5"
)

func NewThreadForm(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sections, err := models.GetAllSectionsWithLastPost()
	if err != nil {
		http.Error(w, "could not load sections", http.StatusInternalServerError)
		return
	}

	Render(w, "new", map[string]any{
		"User":     user,
		"Sections": sections,
	})
}

func CreateThread(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	sectionStr := r.PostForm.Get("section_id")
	sectionID, _ := strconv.Atoi(sectionStr)
	if sectionID == 0 {
		sectionID = 1 // fallback default section
	}

	if title == "" || content == "" {
		http.Error(w, "title and content required", http.StatusBadRequest)
		return
	}

	_, err := models.CreateThread(title, content, user.ID, sectionID)
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

	var user *models.User
	if u := r.Context().Value("user"); u != nil {
		user = u.(*models.User)
	}

	Render(w, "view_thread", map[string]any{
		"Thread":  thread,
		"Replies": replies,
		"User":    user,
	})
}
