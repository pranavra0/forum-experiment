package models

import (
    "crypto/rand"
    "encoding/base64"
    "forum-experiment/db"
    "time"
)

type Session struct {
    ID        int64
    UserID    int64
    Token     string
    CreatedAt time.Time
}

func GenerateToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID int64) (string, error) {
    token, err := GenerateToken()
    if err != nil {
        return "", err
    }

    _, err = db.Conn.Exec(
        "INSERT INTO sessions (user_id, token, created_at) VALUES (?, ?, ?)",
        userID,
        token,
        time.Now().Format(time.RFC3339Nano),
    )
    if err != nil {
        return "", err
    }

    return token, nil
}

func GetUserBySessionToken(token string) (*User, error) {
    var user User
    var createdAtStr string

    err := db.Conn.QueryRow(`
        SELECT u.id, u.username, u.email, u.password_hash, u.created_at 
        FROM users u
        JOIN sessions s ON s.user_id = u.id
        WHERE s.token = ?
    `, token).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAtStr)

    if err != nil {
        return nil, err
    }

    user.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func DeleteSession(token string) error {
    _, err := db.Conn.Exec("DELETE FROM sessions WHERE token = ?", token)
    return err
}