package main

import (
	"fmt"
	"log"
	"forum-experiment/db"
	"forum-experiment/models"
)

func main() {
	db.Init("forum.db")

	const userID = 1

	for i := 1; i <= 25; i++ {
		title := fmt.Sprintf("Test Thread #%d", i)
		content := fmt.Sprintf("This is a sample post for thread #%d.", i)
		_, err := models.CreateThread(title, content, userID)
		if err != nil {
			log.Fatalf("failed to create thread %d: %v", i, err)
		}
	}

	log.Println("✅ Seed complete — 25 threads created.")
}
