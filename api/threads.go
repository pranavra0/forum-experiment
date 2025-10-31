package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"forum-experiment/models"

	"github.com/go-chi/chi/v5"
)

func GetThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := models.GetAllThreads()
	if err != nil {
		http.Error(w, "Failed to get threads", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(threads)
}

func GetThreadByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	thread, err := models.GetThreadByID(id)
	if err != nil || thread == nil {
		http.Error(w, "Thread not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(thread)
}

func DeleteThread(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		log.Println("user not in context")
	} else {
		log.Printf("user found: %s, admin=%v", user.Username, user.IsAdmin)
	}
	if !ok || user == nil || !user.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteThread(id); err != nil {
		http.Error(w, "Failed to delete thread", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
