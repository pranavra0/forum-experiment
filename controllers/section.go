package controllers

import (
	"log"
	"net/http"
	"strconv"

	"forum-experiment/models"

	"github.com/go-chi/chi/v5"
)

func SectionHandler(w http.ResponseWriter, r *http.Request) {

	sectionIDStr := chi.URLParam(r, "id")
	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	user, _ := r.Context().Value("user").(*models.User)

	pageStr := r.URL.Query().Get("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	section, err := models.GetSectionByID(sectionID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	const pageSize = 10
	threads, totalPages, err := models.GetPaginatedThreadsBySection(sectionID, page, pageSize)
	if err != nil {
		http.Error(w, "Error loading threads", http.StatusInternalServerError)
		return
	}

	var ids []int
	for _, t := range threads {
		ids = append(ids, t.ID)
	}

	replyCounts, err := models.GetReplyCountForThreads(ids)
	if err != nil {
		log.Printf("⚠️ Error fetching reply counts: %v", err)
	} else {
		for i := range threads {
			if count, ok := replyCounts[threads[i].ID]; ok {
				threads[i].ReplyCount = count
			}
		}
	}

	pagination := BuildPagination(page, totalPages)

	data := map[string]any{
		"User":              user,
		"Section":           section,
		"Threads":           threads,
		"Page":              pagination.Page,
		"TotalPages":        pagination.TotalPages,
		"Pages":             pagination.Pages,
		"ShowStartEllipsis": pagination.ShowStartEllipsis,
		"ShowEndEllipsis":   pagination.ShowEndEllipsis,
		"HasPrev":           pagination.HasPrev,
		"HasNext":           pagination.HasNext,
	}

	Render(w, "section", data)
}
