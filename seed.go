package main

// right now this file is just for testing stuff
import (
	"fmt"
	"log"

	"forum-experiment/db"
	"forum-experiment/models"
)

func main() {
	if err := db.Init("forum.db"); err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer db.Close()

	log.Println("🌱 Seeding database...")

	// Sections
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
		_, err := db.Conn.Exec(`INSERT OR IGNORE INTO sections (name, description) VALUES (?, ?)`, s.Name, s.Description)
		if err != nil {
			log.Fatalf("failed to create section %q: %v", s.Name, err)
		}
	}
	log.Println("Sections created.")

	// Create admin user (username: admin, password: abc123)
	adminUser, err := models.GetUserByUsername("admin")
	if err != nil {
		log.Fatalf("failed checking admin user: %v", err)
	}
	if adminUser == nil {
		if err := models.CreateUser("admin", "admin@example.com", "abc123", true); err != nil {
			log.Fatalf("failed to create admin user: %v", err)
		}
		log.Println("✅ Created admin user: username=admin, password=abc123")
	} else {
		log.Println("🔒 Admin user already exists; skipping creation.")
	}

	// Test user (username: testuser)
	testUser, err := models.GetUserByUsername("testuser")
	if err != nil {
		log.Fatalf("failed checking test user: %v", err)
	}
	if testUser == nil {
		if err := models.CreateUser("testuser", "test@example.com", "testpass", false); err != nil {
			log.Fatalf("failed to create test user: %v", err)
		}
		log.Println("✅ Created test user: username=testuser, password=testpass")
	} else {
		log.Println("🔒 Test user already exists; skipping creation.")
	}

	// Fetch userID for testuser
	var userID int64
	u, err := models.GetUserByUsername("testuser")
	if err != nil || u == nil {
		log.Fatalf("failed to fetch test user id: %v", err)
	}
	userID = u.ID

	// Fetch sections
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
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			log.Fatalf("failed scanning section row: %v", err)
		}
		sectionList = append(sectionList, s)
	}

	// Create threads and replies
	for _, s := range sectionList {
		if s.Name == "Books" {
			log.Printf("📖 Leaving section %q empty for testing.", s.Name)
			continue
		}

		// Create threads
		for i := 1; i <= 8; i++ {
			title := fmt.Sprintf("[%s] Sample Thread #%d", s.Name, i)
			content := fmt.Sprintf("Discussion topic #%d in the %s section.", i, s.Name)
			threadID, err := models.CreateThread(title, content, userID, s.ID)
			if err != nil {
				log.Fatalf("failed to create thread in section %q: %v", s.Name, err)
			}

			// Create root replies
			for r := 1; r <= 3; r++ {
				rootContent := fmt.Sprintf("Root reply #%d for thread %d", r, threadID)
				if err := models.CreateReply(int(threadID), userID, rootContent, nil); err != nil {
					log.Fatalf("failed to create root reply for thread %d: %v", threadID, err)
				}
			}

			// Fetch root replies for children
			rootReplies, err := models.GetRepliesByThreadID(int(threadID))
			if err != nil {
				log.Fatalf("failed to fetch replies for thread %d: %v", threadID, err)
			}

			for _, r := range rootReplies {
				// Add 2 child replies per root
				for c := 1; c <= 2; c++ {
					childContent := fmt.Sprintf("Child reply #%d to reply %d", c, r.ID)
					if err := models.CreateReply(int(threadID), userID, childContent, &r.ID); err != nil {
						log.Fatalf("failed to create child reply for reply %d: %v", r.ID, err)
					}
				}
			}
		}
		log.Printf("✅ Created 8 threads with nested replies in section %q.", s.Name)
	}

	log.Println("Seeding complete!")
}
