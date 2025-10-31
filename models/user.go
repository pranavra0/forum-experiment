package models

import (
	"database/sql"
	"time"

	"forum-experiment/db"

	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
}

func CreateUser(username, email, password string, isAdmin bool) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(
		`INSERT INTO users (username, email, password_hash, is_admin, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		username,
		email,
		string(hash),
		isAdmin,
		time.Now().Format(time.RFC3339Nano),
	)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	var u User
	var createdAtStr string
	var isAdminInt int

	err := db.Conn.QueryRow(
		`SELECT id, username, email, password_hash, is_admin, created_at
		 FROM users WHERE username = ?`, username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &isAdminInt, &createdAtStr)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	u.IsAdmin = isAdminInt != 0 // convert int -> bool

	u.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func CheckPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

func EnsureAdminExists(username, email, password string) error {
	u, err := GetUserByUsername(username)
	if err != nil {
		return err
	}
	if u != nil {
		log.Printf("Admin user %s already exists.", username)
		return nil
	}

	log.Printf("Creating admin user: %s", username)
	return CreateUser(username, email, password, true)
}
