package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"forum-experiment/api"
	"forum-experiment/config"
	"forum-experiment/controllers"
	"forum-experiment/db"
	"forum-experiment/models"
)

func main() {

	config.LoadConfig()
	config.LoadEnv()

	dbPath := config.E.DatabasePath
	if config.C.DB.Path != "" && config.E.DatabasePath == "forum.db" {
		dbPath = config.C.DB.Path
	}

	if err := db.Init(dbPath); err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer db.Close()

	if err := models.EnsureAdminExists(
		config.E.AdminUsername,
		config.E.AdminEmail,
		config.E.AdminPassword,
	); err != nil {
		log.Fatalf("failed to ensure admin user: %v", err)
	}

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

	port := config.E.Port
	if port == "" {
		port = config.C.App.Port // fallback if YAML specifies one
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Server running at %s", addr)
	http.ListenAndServe(addr, r)
}
