package main

// file for testing sections + pagination

import (
	"fmt"
	"log"
	"time"

	"forum-experiment/db"
	"forum-experiment/models"
)

func main() {
	if err := db.Init("forum.db"); err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer db.Close()

	log.Println("🌱 Seeding database...")

	sections := []struct {
		Name        string
		Description string
	}{
		{"Sports", "I just spiked my soda like im in NFL"},
		{"Life", "Yeah"},
		{"Tech", "Tech SXN"},
		{"Books", "Books are banned"}, // will stay empty
	}

	for _, s := range sections {
		_, err := db.Conn.Exec(`INSERT INTO sections (name, description) VALUES (?, ?)`, s.Name, s.Description)
		if err != nil {
			log.Fatalf("failed to create section %q: %v", s.Name, err)
		}
	}
	log.Println("Sections created.")

	_, err := db.Conn.Exec(`
		INSERT OR IGNORE INTO users (username, email, password_hash, created_at)
		VALUES ('testuser', 'test@example.com', 'hashedpassword', ?)
	`, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		log.Fatalf("failed to create test user: %v", err)
	}

	var userID int
	if err := db.Conn.QueryRow(`SELECT id FROM users WHERE username = 'testuser'`).Scan(&userID); err != nil {
		log.Fatalf("failed to fetch test user id: %v", err)
	}

	rows, err := db.Conn.Query(`SELECT id, name FROM sections`)
	if err != nil {
		log.Fatalf("failed to fetch section IDs: %v", err)
	}
	defer rows.Close()

	type sec struct {
		ID   int
		Name string
	}
	var sectionList []sec
	for rows.Next() {
		var s sec
		rows.Scan(&s.ID, &s.Name)
		sectionList = append(sectionList, s)
	}

	for _, s := range sectionList {
		if s.Name == "Books" {
			log.Printf("📖 Leaving section %q empty for testing.", s.Name)
			continue
		}

		for i := 1; i <= 8; i++ {
			title := fmt.Sprintf("[%s] Sample Thread #%d", s.Name, i)
			content := fmt.Sprintf("Discussion topic #%d in the %s section.", i, s.Name)
			_, err := models.CreateThread(title, content, int64(userID), s.ID)
			if err != nil {
				log.Fatalf("failed to create thread in section %q: %v", s.Name, err)
			}
		}
		log.Printf("✅ Created 8 threads in section %q.", s.Name)
	}

	log.Println("Seeding complete!")
}
