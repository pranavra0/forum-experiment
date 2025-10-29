package models

import (
	"database/sql"
	"fmt"
	"forum-experiment/db"
	"time"
)

type Reply struct {
	ID        int
	ThreadID  int
	UserID    int
	Username  string
	ParentID  *int // nil if root-level reply
	Content   string
	CreatedAt time.Time
	Children  []*Reply // populated in hierarchical fetch
}

func CreateReply(threadID int, userID int64, content string, parentID *int) error {
	_, err := db.Conn.Exec(
		"INSERT INTO replies (thread_id, user_id, parent_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		threadID, userID, parentID, content, time.Now().Format(time.RFC3339Nano),
	)
	return err
}

func GetRepliesByThreadID(threadID int) ([]*Reply, error) {
	rows, err := db.Conn.Query(`
        SELECT r.id, r.thread_id, r.user_id, r.parent_id, u.username, r.content, r.created_at
        FROM replies r
        JOIN users u ON r.user_id = u.id
        WHERE r.thread_id = ?
        ORDER BY r.created_at ASC
    `, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allReplies := []*Reply{}
	replyMap := map[int]*Reply{}

	for rows.Next() {
		var r Reply
		var parentID sql.NullInt64
		var created string
		if err := rows.Scan(&r.ID, &r.ThreadID, &r.UserID, &parentID, &r.Username, &r.Content, &created); err != nil {
			return nil, err
		}
		if parentID.Valid {
			val := int(parentID.Int64)
			r.ParentID = &val
		}
		r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)

		// Attach quote if parent exists
		if r.ParentID != nil {
			if parent, ok := replyMap[*r.ParentID]; ok {
				r.Content = fmt.Sprintf("[quote=%s]%s[/quote]\n%s", parent.Username, parent.Content, r.Content)
			}
		}

		allReplies = append(allReplies, &r)
		replyMap[r.ID] = &r
	}

	return allReplies, nil
}

func (r Reply) FormattedTime() string {
	return r.CreatedAt.Format("2006-01-02 15:04")
}

func GetPaginatedRepliesByThreadID(threadID, page, pageSize int) ([]*Reply, int, error) {
	allReplies, err := GetRepliesByThreadID(threadID)
	if err != nil {
		return nil, 0, err
	}
	total := len(allReplies)
	totalPages := (total + pageSize - 1) / pageSize

	start := (page - 1) * pageSize
	if start >= total {
		return []*Reply{}, totalPages, nil
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	paginated := allReplies[start:end]
	return paginated, totalPages, nil
}
