package models

import (
	"time"
	"forum-experiment/db"
)

type Reply struct {
	ID        int
	ThreadID  int
	Content   string
	CreatedAt time.Time
}

func CreateReply(threadID int, content string) error {
	_, err := db.Conn.Exec(
		"INSERT INTO replies (thread_id, content, created_at) VALUES (?, ?, ?)",
		threadID, content, time.Now().Format(time.RFC3339Nano),
	)
	return err
}

func GetRepliesByThreadID(threadID int) ([]Reply, error) {
	rows, err := db.Conn.Query(
		"SELECT id, thread_id, content, created_at FROM replies WHERE thread_id = ? ORDER BY created_at ASC",
		threadID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var r Reply
		var created string
		if err := rows.Scan(&r.ID, &r.ThreadID, &r.Content, &created); err != nil {
			return nil, err
		}
		r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		replies = append(replies, r)
	}
	return replies, rows.Err()
}
