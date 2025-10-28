package controllers

import (
	"net/http"

	"forum-experiment/models"
)

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
	user, _ := r.Context().Value("user").(*models.User)

	sections, err := models.GetAllSectionsWithLastPost()
	if err != nil {
		http.Error(w, "Failed to load sections", http.StatusInternalServerError)
		return
	}

	Render(w, "home", map[string]any{
		"User":     user,
		"Sections": sections,
	})
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*models.User)
	query := r.URL.Query().Get("q")

	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	results, err := models.SearchThreads(query)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	Render(w, "search_results", map[string]any{
		"User":    user,
		"Query":   query,
		"Results": results,
	})
}
