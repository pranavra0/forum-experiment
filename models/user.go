package models

import (
	"database/sql"
	"time"
	"golang.org/x/crypto/bcrypt"
	"forum-experiment/db"
)

type User struct {
    ID           int64
    Username     string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
}

func CreateUser(username, email, password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = db.Conn.Exec(
        "INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
        username,
        email,
        string(hash),
        time.Now().Format(time.RFC3339Nano),
    )
    return err
}

func GetUserByUsername(username string) (*User, error) {
    var u User
    var createdAtStr string
    
    err := db.Conn.QueryRow(
        "SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?", 
        username,
    ).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &createdAtStr)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

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