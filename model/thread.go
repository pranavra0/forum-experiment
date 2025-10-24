package model

import (
	"time"
	"log"

	"forum-experiment/db"
)

type Thread struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

func GetAllThreads() ([]Thread, error) {
	rows, err := db.Conn.Query("SELECT id, title, content, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []Thread
	for rows.Next() {
		var t Thread
		var created string
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &created); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		threads = append(threads, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return threads, nil
}

func CreateThread(title, content string) (int64, error) {
	stmt, err := db.Conn.Prepare("INSERT INTO threads (title, content, created_at) VALUES (?, ?, ?)")
	log.Printf("✅ New thread created: %s — %s\n", title, content)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(title, content, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
