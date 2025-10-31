package controllers

import (
	"context"
	"forum-experiment/models"
	"log"
	"net/http"
)

const sessionCookie = "session"

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil {
			cookie, err = r.Cookie("session_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		user, err := models.GetUserBySessionToken(cookie.Value)
		if err != nil || user == nil {
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("Found user: %s", user.Username)
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
