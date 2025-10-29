package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // no cgo lol
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
		user_id INTEGER NOT NULL,
		section_id INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS replies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		thread_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		parent_id INTEGER DEFAULT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES replies(id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS sections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT
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
