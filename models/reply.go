package models

import (
	"forum-experiment/db"
	"time"
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

func (r Reply) FormattedTime() string {
	return r.CreatedAt.Format("2006-01-02 15:04")
}

func GetPaginatedRepliesByThreadID(threadID, page, pageSize int) ([]Reply, int, error) {
	offset := (page - 1) * pageSize

	rows, err := db.Conn.Query(`
		SELECT r.id, r.thread_id, r.user_id, u.username, r.content, r.created_at
		FROM replies r
		JOIN users u ON r.user_id = u.id
		WHERE r.thread_id = ?
		ORDER BY r.created_at ASC
		LIMIT ? OFFSET ?
	`, threadID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var r Reply
		var created string
		if err := rows.Scan(&r.ID, &r.ThreadID, &r.UserID, &r.Username, &r.Content, &created); err != nil {
			return nil, 0, err
		}
		r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		replies = append(replies, r)
	}

	var total int
	if err := db.Conn.QueryRow(`SELECT COUNT(*) FROM replies WHERE thread_id = ?`, threadID).Scan(&total); err != nil {
		return nil, 0, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	return replies, totalPages, nil
}
