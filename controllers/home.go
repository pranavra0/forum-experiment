package controllers

import (
	"log"
	"net/http"
	"strconv"

	"forum-experiment/models"
)

// --- Pagination logic ---
func BuildPagination(page, totalPages int) Pagination {
	if totalPages <= 0 {
		return Pagination{}
	}

	var pages []int
	const window = 1

	start := page - window
	end := page + window

	if start < 2 {
		start = 2
	}
	if end > totalPages-1 {
		end = totalPages - 1
	}

	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}

	showStartEllipsis := len(pages) > 0 && pages[0] > 2
	showEndEllipsis := len(pages) > 0 && pages[len(pages)-1] < totalPages-1

	return Pagination{
		Page:              page,
		TotalPages:        totalPages,
		Pages:             pages,
		ShowStartEllipsis: showStartEllipsis,
		ShowEndEllipsis:   showEndEllipsis,
		HasPrev:           page > 1,
		HasNext:           page < totalPages,
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	user, _ := r.Context().Value("user").(*models.User)
	const pageSize = 10

	threads, totalPages, err := models.GetPaginatedThreads(page, pageSize)
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

	log.Printf("DEBUG pagination: page=%d totalPages=%d pages=%v", page, totalPages, pagination.Pages)

	data := map[string]any{
		"User":              user,
		"Threads":           threads,
		"Page":              pagination.Page,
		"TotalPages":        pagination.TotalPages,
		"Pages":             pagination.Pages,
		"ShowStartEllipsis": pagination.ShowStartEllipsis,
		"ShowEndEllipsis":   pagination.ShowEndEllipsis,
		"HasPrev":           pagination.HasPrev,
		"HasNext":           pagination.HasNext,
	}

	Render(w, "home", data)
}
