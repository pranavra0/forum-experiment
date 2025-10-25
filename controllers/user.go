package controllers

import (
    "net/http"
    "forum-experiment/models"
	"log"
)

func ShowRegister(w http.ResponseWriter, r *http.Request) {
    Render(w, "register", PageData{
        Name: "register",
    })
}

func Register(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    // Validation
    if username == "" || email == "" || password == "" {
        Render(w, "register", PageData{
            Name:  "register",
            Error: "All fields are required",
        })
        return
    }

    // Check if username already exists
    existingUser, err := models.GetUserByUsername(username)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    if existingUser != nil {
        Render(w, "register", PageData{
            Name:  "register",
            Error: "Username already taken",
        })
        return
    }

    // Create user
    err = models.CreateUser(username, email, password)
    if err != nil {
        Render(w, "register", PageData{
            Name:  "register",
            Error: "Could not create user",
        })
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
    Render(w, "login", PageData{
        Name: "login",
    })
}

func Login(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    if username == "" || password == "" {
        Render(w, "login", PageData{
            Name:  "login",
            Error: "Username and password are required",
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
        Render(w, "login", PageData{
            Name:  "login",
            Error: "Invalid username or password",
        })
        return
    }

    // Create session
    token, err := models.CreateSession(user.ID)
    if err != nil {
        log.Printf("Error creating session: %v", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Set cookie
    http.SetCookie(w, &http.Cookie{
        Name:     sessionCookie,
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        MaxAge:   86400 * 30, // 30 days
    })

    log.Printf("Login successful for user: %s", user.Username)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func Logout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie(sessionCookie)
    if err == nil {
        models.DeleteSession(cookie.Value)
    }
    
    http.SetCookie(w, &http.Cookie{
        Name:     sessionCookie,
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
    })
    
    http.Redirect(w, r, "/", http.StatusSeeOther)
}