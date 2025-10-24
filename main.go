package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"forum-experiment/controller"
	"forum-experiment/db"
)

func main() {
	if err := db.Init("forum.db"); err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	// Routes
	r.Get("/", controllers.Home)
	r.Get("/thread/new", controllers.NewThreadForm)
	r.Post("/thread/new", controllers.CreateThread)

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
