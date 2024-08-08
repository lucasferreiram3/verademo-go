package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"verademo-go/src-app/shared/view"
)

type User struct {
	Username     string
	PasswordHint string
	CreatedAt    time.Time
	LastLogin    time.Time
	RealName     string
	BlabName     string
}

var db *sql.DB

func getMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
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
			http.Redirect(w, req, "feed.html", http.StatusFound)
		}
		return
	}

	user, err := createFromRequest(req)

	if err != nil {
		log.Println("Error creating user from request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user != nil {
		http.SetCookie(w, &http.Cookie{Name: "session_username", Value: user.Username})
		log.Println("User is remembered - redirecting...")
		if target != "" {
			http.Redirect(w, req, target, http.StatusFound)
		} else {
			http.Redirect(w, req, "feed.html", http.StatusFound)
		}
		return
	} else {
		log.Println("User is not remembered")
	}

	if username == "" {
		username = ""
	}

	if target == "" {
		target = ""
	}
	log.Println("Entering showLogin with username %s and target %s\n", username, target)

	view.Render(w, "login.html", []byte{})
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

	// Constructing SQL Query, have to figure out hashing
	sqlQuery := fmt.Sprintf("SELECT username, password, password_hint, created_at, last_login, real_name, blab_name FROM users WHERE username='%s' AND password='%s';", username, getMD5(password))
	log.Println("Executing SQL Query:", sqlQuery)

	rows, err := db.Query(sqlQuery)

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

		if remember != "" {
			// updateInResponse needs implementation
		}

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

func processLogout(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering processLogout")

	// Clear session username
	http.SetCookie(w, &http.Cookie{Name: "session_username", Value: "", Path: "/", MaxAge: -1})

	currentUser := &User{}

	if err := updateInResponse(currentUser, w); err != nil {
		log.Println("Error updating response:", err)
	}

	http.Redirect(w, req, "/login", http.StatusFound)

}

func showPasswordHint(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	log.Println("Entering password-hint with username: " + username)

	if username != "" {
		return
	}

}

func showRegister(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering showRegister")

	view.Render(w, "register.html", []byte{})
}

func showRegisterFinish(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering showRegisterFinish")

	http.Redirect(w, req, "/register-finish", http.StatusFound)
}

func processRegisterFinish(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering processRegisterFinish")

}

func createFromRequest(req *http.Request) (*User, error) {
	cookie, err := req.Cookie("user")
	if err != nil {
		log.Println("No user cookie.")
		return nil, nil
	}

	cookieValue := cookie.Value
	decoded, err := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		log.Println("Error decoding cookie:", err)
		return nil, err
	}

	var user User
	if err := json.Unmarshal(decoded, &user); err != nil {
		log.Println("Error unmarshaling user from cookie:", err)
		return nil, err
	}

	log.Println("Username is:", user.Username)
	return &user, nil
}

func updateInResponse(currentUser *User, w http.ResponseWriter) error {
	userJSON, err := json.Marshal(currentUser)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(userJSON)
	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: encoded,
		Path:  "/",
	})

	return nil
}
