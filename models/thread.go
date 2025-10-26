package models

import (
	"time"
	"log"
	"forum-experiment/db"
	"math"
)

type Thread struct {
	ID        int
	Title     string
	Content   string
	UserID    int
	Username  string 
	CreatedAt time.Time
}

func GetAllThreads() ([]Thread, error) {
	rows, err := db.Conn.Query(`
		SELECT t.id, t.title, t.content, t.created_at, t.user_id, u.username
		FROM threads t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []Thread
	for rows.Next() {
		var t Thread
		var created string
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &created, &t.UserID, &t.Username); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		threads = append(threads, t)
	}
	return threads, rows.Err()
}

func CreateThread(title, content string, userID int64) (int64, error) {
	stmt, err := db.Conn.Prepare("INSERT INTO threads (title, content, user_id, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(title, content, userID, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}

	log.Printf("✅ User %d created thread: %s", userID, title)
	return res.LastInsertId()
}


func GetThreadByID(id int) (Thread, error) {
	var t Thread
	var created string
	err := db.Conn.QueryRow(`
		SELECT t.id, t.title, t.content, t.created_at, t.user_id, u.username
		FROM threads t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`, id).Scan(&t.ID, &t.Title, &t.Content, &created, &t.UserID, &t.Username)
	if err == nil {
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	}
	return t, err
}

func (t Thread) FormattedTime() string {
    return t.CreatedAt.Format("2006-01-02 15:04")
}

// pagination logic

func GetThreadsPage(limit, offset int) ([]Thread, error) {
	rows, err := db.Conn.Query(`
		SELECT t.id, t.title, t.content, t.created_at, t.user_id, u.username
		FROM threads t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []Thread
	for rows.Next() {
		var t Thread
		var created string
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &created, &t.UserID, &t.Username); err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		threads = append(threads, t)
	}
	return threads, rows.Err()
}

func CountThreads() (int, error) {
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM threads").Scan(&count)
	return count, err
}

func GetPaginatedThreads(page, pageSize int) ([]Thread, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	total, err := CountThreads()
	if err != nil {
		return nil, 0, err
	}

	// Fetch threads for this page
	rows, err := db.Conn.Query(`
		SELECT t.id, t.title, t.content, t.created_at, t.user_id, u.username
		FROM threads t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var threads []Thread
	for rows.Next() {
		var t Thread
		var created string
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &created, &t.UserID, &t.Username); err != nil {
			return nil, 0, err
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		threads = append(threads, t)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return threads, totalPages, nil
}

