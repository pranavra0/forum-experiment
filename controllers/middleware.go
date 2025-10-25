package controllers

import (
    "context"
    "forum-experiment/models"
    "net/http"
	"log"
)

const sessionCookie = "session"

func WithUser(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie(sessionCookie)
        if err != nil {
            next.ServeHTTP(w, r)
            return
        }

        user, err := models.GetUserBySessionToken(cookie.Value)
        if err != nil {
            next.ServeHTTP(w, r)
            return
        }
        log.Printf("Found user: %s", user.Username)

        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}