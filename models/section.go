package models

import (
	"forum-experiment/db"
	"log"
	"time"
)

type PostSummary struct {
	ID        int
	Title     string
	Username  string
	CreatedAt string
}

type Section struct {
	ID          int
	Name        string
	Description string
	LastPost    *PostSummary // pointer so it can be nil when no posts exist
}

func GetAllSectionsWithLastPost() ([]Section, error) {
	rows, err := db.Conn.Query(`
		SELECT
			s.id,
			s.name,
			IFNULL(s.description, ''),
			r.id,              -- reply ID
			t.id,              -- thread ID
			t.title,           -- thread title
			u.username,        -- reply author
			r.created_at       -- reply creation time
		FROM sections s
		LEFT JOIN replies r ON r.id = (
			SELECT rr.id
			FROM replies rr
			JOIN threads tt ON tt.id = rr.thread_id
			WHERE tt.section_id = s.id
			ORDER BY rr.created_at DESC
			LIMIT 1
		)
		LEFT JOIN threads t ON t.id = r.thread_id
		LEFT JOIN users u ON u.id = r.user_id
		ORDER BY s.name ASC
	`)
	if err != nil {
		log.Printf("❌ GetAllSectionsWithLastPost SQL error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var s Section
		var replyID, threadID *int
		var threadTitle, username, createdAt *string

		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &replyID, &threadID, &threadTitle, &username, &createdAt); err != nil {
			return nil, err
		}

		if replyID != nil {
			s.LastPost = &PostSummary{
				ID:        *threadID,
				Title:     *threadTitle,
				Username:  *username,
				CreatedAt: *createdAt,
			}
		}

		sections = append(sections, s)
	}

	return sections, rows.Err()
}

// helper to get section by id
func GetSectionByID(id int) (Section, error) {
	var s Section
	err := db.Conn.QueryRow(`SELECT id, name, IFNULL(description, '') FROM sections WHERE id = ?`, id).
		Scan(&s.ID, &s.Name, &s.Description)
	return s, err
}

func (p *PostSummary) FormattedTime() string {
	if p == nil || p.CreatedAt == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339Nano, p.CreatedAt)
	if err != nil {
		return p.CreatedAt
	}
	return parsed.Format("2006-01-02 15:04")
}
