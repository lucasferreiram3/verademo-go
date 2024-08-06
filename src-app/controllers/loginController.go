package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username     string
	PasswordHint string
	CreatedAt    time.Time
	LastLogin    time.Time
	RealName     string
	BlabName     string
}

func showLogin(w http.ResponseWriter, req *http.Request) {
	target := req.URL.Query().Get("target")
	username := req.URL.Query().Get("username")

	session, err := req.Cookie("session_username")
	if err == nil && session.Value != "" {
		log.Println("User is already logged in - redirecting...")
		if target != "" {
			http.Redirect(w, req, target, http.StatusFound)
		} else {
			http.Redirect(w, req, "/feed", http.StatusFound)
		}
		return
	}
	// TODO: Code below written is based off the JS version, need to see if theres a createRequest equivalent in GO
	// user, err := createFromRequest(req)

	if err != nil {
		log.Println("Error creating user from request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var user User
	if user != nil {
		http.SetCookie(w, &http.Cookie{Name: "session_username", Value: user.Username})
		log.Println("User is remembered - redirecting...")
		if target != "" {
			http.Redirect(w, req, target, http.StatusFound)
		} else {
			http.Redirect(w, req, "/feed", http.StatusFound)
		}
		return
	} else {
		log.Println("User is not remembered")
	}

	if username == "" {
		username == ""
	}

	if target == "" {
		target == ""
	}
	log.Println("Entering showLogin with username %s and target %s\n", username, target)
}

func processLogin(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering processLogin")

	// Form data check
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	remember := req.FormValue("remember")
	target := req.FormValue("target")

	var nextView string
	if target != "" {
		nextView = target
	} else {
		nextView = "/feed"
	}

	// Constructing SQL Query
	sqlQuery := fmt.Sprintf("SELECT username, password, password_hint, created_at, last_login, real_name, blab_name FROM users WHERE username='%s' AND password='%s';")
	log.Println("Executing SQL Query:", sqlQuery)

	// Execute SQL Query
	if err != nil {
		log.Println("Error executing SQL Query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Check if user is found
	var user User
	if rows.Next() {
		if err := rows.Scan(&user.Username, &user.PasswordHint, &user.CreatedAt, &user.LastLogin, &user.RealName, &user.BlabName); err != nil {
			log.Println("Error scanning result:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Println("User found:", user.Username)
		http.SetCookie(w, &http.Cookie{Name: "username", Value: user.Username})

		if len(username) >= 4 && username[len(username)-4:] == "totp" {
			log.Println("User has TOTP enabled")
			req.Header.Set("totp_username", user.Username)
			nextView = "/totp"
		} else {
			log.Println("Setting session username to:", user.Username)
			// Update last login
			_, err = db.Exec("UPDATE users SET last_login=NOW() WHERE username=?", user.Username)
			if err != nil {
				log.Println("Error updating last login:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// Set session username
			http.SetCookie(w, &http.Cookie{Name: "session_username", Value: user.Username})
		}
	} else {
		log.Println("User Not Found")
		http.Error(w, "Login failed. Please try again.", http.StatusUnauthorized)
		nextView = "/login"
	}
	log.Println("Redirecting to view:", nextView)
	http.Redirect(w, req, nextView, http.StatusFound)
}
