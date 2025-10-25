package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"forum-experiment/models"
)

func CreateReply(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/thread/")
	idStr = strings.TrimSuffix(idStr, "/reply")
	threadID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid thread id", http.StatusBadRequest)
		return
	}

	content := r.PostForm.Get("content")
	if content == "" {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}

	if err := models.CreateReply(threadID, user.ID, content); err != nil {
		http.Error(w, "could not save reply", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/thread/"+strconv.Itoa(threadID), http.StatusSeeOther)
}
