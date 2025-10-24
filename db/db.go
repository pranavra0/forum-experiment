package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite" // no cgo lol
	"time"
)

var Conn *sql.DB

func Init(path string) error {
	var err error
	Conn, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	Conn.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS threads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	`
	if _, err := Conn.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	if err := Conn.Ping(); err != nil {
		return err
	}

	Conn.Exec("PRAGMA busy_timeout = 5000;")
	_ = time.Now
	return nil
}

func Close() error {
	if Conn != nil {
		return Conn.Close()
	}
	return nil
}
