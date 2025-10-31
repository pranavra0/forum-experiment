package controllers

import (
	"log"
	"net/http"

	"forum-experiment/models"
)

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	Render(w, "register", map[string]any{})
}

// Register handles new user creation
func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		Render(w, "register", map[string]any{
			"Error": "All fields are required",
		})
		return
	}

	existingUser, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		Render(w, "register", map[string]any{
			"Error": "Username already taken",
		})
		return
	}

	err = models.CreateUser(username, email, password, false)
	if err != nil {
		Render(w, "register", map[string]any{
			"Error": "Could not create user",
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	Render(w, "login", map[string]any{})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		Render(w, "login", map[string]any{
			"Error": "Username and password are required",
		})
		return
	}

	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if user == nil || !models.CheckPassword(user, password) {
		Render(w, "login", map[string]any{
			"Error": "Invalid username or password",
		})
		return
	}

	token, err := models.CreateSession(user.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookie)
	if err == nil {
		models.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   sessionCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
