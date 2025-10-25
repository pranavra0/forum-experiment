package models

import (
	"time"
	"forum-experiment/db"
)

type Reply struct {
	ID        int
	ThreadID  int
	UserID    int
	Username  string
	Content   string
	CreatedAt time.Time
}

func CreateReply(threadID int, userID int64, content string) error {
	_, err := db.Conn.Exec(
		"INSERT INTO replies (thread_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		threadID, userID, content, time.Now().Format(time.RFC3339Nano),
	)
	return err
}


func GetRepliesByThreadID(threadID int) ([]Reply, error) {
	rows, err := db.Conn.Query(`
		SELECT r.id, r.thread_id, r.user_id, u.username, r.content, r.created_at
		FROM replies r
		JOIN users u ON r.user_id = u.id
		WHERE r.thread_id = ?
		ORDER BY r.created_at ASC
	`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var r Reply
		var created string
		if err := rows.Scan(&r.ID, &r.ThreadID, &r.UserID, &r.Username, &r.Content, &created); err != nil {
			return nil, err
		}
		r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		replies = append(replies, r)
	}
	return replies, rows.Err()
}
