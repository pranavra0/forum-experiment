package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"forum-experiment/api"
	"forum-experiment/controllers"
	"forum-experiment/db"
)

func main() {
	if err := db.Init("forum.db"); err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer db.Close()
	controllers.InitTemplates()
	r := chi.NewRouter()

	// middleware
	r.Use(controllers.WithUser)

	// Routes
	r.Get("/", controllers.HomeHandler)
	r.Get("/section/{id}", controllers.SectionHandler)
	r.Get("/thread/new", controllers.NewThreadForm)
	r.Post("/thread/new", controllers.CreateThread)
	r.Get("/thread/{id}", controllers.ShowThread)
	r.Post("/thread/{id}/reply", controllers.CreateReply)
	r.Get("/register", controllers.ShowRegister)
	r.Post("/register", controllers.Register)
	r.Get("/login", controllers.ShowLogin)
	r.Post("/login", controllers.Login)
	r.Get("/logout", controllers.Logout)
	r.Post("/logout", controllers.Logout)
	r.Get("/search", controllers.SearchHandler)

	//api stuff
	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Post("/login", api.Login)
		apiRouter.Get("/threads", api.GetThreads)
		apiRouter.Get("/threads/{id}", api.GetThreadByID)
		apiRouter.Delete("/threads/{id}", api.DeleteThread)
	})

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
