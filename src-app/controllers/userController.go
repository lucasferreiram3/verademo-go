package controllers

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"verademo-go/src-app/models"
	sqlite "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	view "verademo-go/src-app/shared/view"

	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

type User struct {
	Username     string
	PasswordHint string
	CreatedAt    time.Time
	LastLogin    time.Time
	RealName     string
	BlabName     string
}
type Account struct {
	Error    string
	Username string
}
type Output struct {
	Username   string
	Error      string
	Hecklers   []models.Blabber
	Events     []string
	Image      string // ImageName
	RealName   string
	BlabName   string
	TotpSecret string
}

func ShowLogin(w http.ResponseWriter, req *http.Request) {
	target := req.URL.Query().Get("target")
	username := req.URL.Query().Get("username")

	current_session, err := req.Cookie("session_username")
	if err == nil && current_session.Value != "" {
		log.Println("User is already logged in - redirecting...")
		if target != "" {
			http.Redirect(w, req, target, http.StatusFound)
		} else {
			http.Redirect(w, req, "/feed", http.StatusFound)
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
			http.Redirect(w, req, "/feed", http.StatusFound)
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
	log.Println("Entering showLogin with username " + username + " and target " + target)

	view.Render(w, "login.html", nil)
}

func ProcessLogin(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering processLogin")

	// Form data check
	err := req.ParseForm()
	if err != nil {
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
	log.Println("Username: " + username + " Password: " + password)
	// Constructing SQL Query
	sqlQuery := "SELECT username, password_hint, created_at, last_login, real_name, blab_name FROM users WHERE username = ? AND password = ?"

	result := struct {
		Username     string
		PasswordHint string
		CreatedAt    string
		LastLogin    string
		RealName     string
		BlabName     string
	}{}

	err = sqlite.DB.QueryRow(sqlQuery, username, GetMD5Hash(password)).Scan(
		&result.Username,
		&result.PasswordHint,
		&result.CreatedAt,
		&result.LastLogin,
		&result.RealName,
		&result.BlabName,
	)
	log.Println("After Query: " + result.Username)
	// In case user does not exist
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found")
			http.Error(w, "Login failed. Please try again.", http.StatusUnauthorized)
			return
		}
		log.Println(err)
		http.Error(w, "An error has occured", http.StatusInternalServerError)
	}
	log.Println("User found. Redirecting...")

	http.SetCookie(w, &http.Cookie{Name: "username", Value: result.Username})

	// Handling the "remember me"
	if remember == "" {
		// Store details in session
		current_session, _ := session.Store.Get(req, "session-name")
		current_session.Values["username"] = result.Username
		current_session.Values["password_hint"] = result.PasswordHint
		current_session.Values["created_at"] = result.CreatedAt
		current_session.Values["last_login"] = result.LastLogin
		current_session.Values["real_name"] = result.RealName
		current_session.Values["blab_name"] = result.BlabName
		current_session.Save(req, w)

		if err := current_session.Save(req, w); err != nil {
			log.Println("Error saving session:", err)
			http.Error(w, "An error occurred", http.StatusInternalServerError)
			return
		}

	}
	// Updating last login time
	_, err = sqlite.DB.Exec("UPDATE users SET last_login=datetime('now') WHERE username='" + result.Username + "';")
	if err != nil {
		log.Println("Error updating last login for user: ", err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	// TOTP Handling
	if len(username) >= 4 && username[len(username)-4:] == "totp" {
		log.Println("User " + username + " Has TOTP Enabled!")
		current_session, _ := session.Store.Get(req, "session-name")
		current_session.Values["totp_username"] = result.Username
		current_session.Save(req, w)
		nextView = "/totp"
	} else {
		log.Println("Setting session username to: " + username)
		current_session, _ := session.Store.Get(req, "session-name")
		current_session.Values["username"] = result.Username
		current_session.Save(req, w)
		nextView = "/feed"
	}

	log.Println("Redirecting to view: " + nextView)
	http.Redirect(w, req, nextView, http.StatusSeeOther)

}

func processLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering processLogout")

	current_session, _ := session.Store.Get(r, "session-name")

	current_session.Values["username"] = ""

	err := current_session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	// Optionally update response
	if err := updateInResponse(current_session.Values["username"], w); err != nil {
		log.Println("Error updating response:", err)
	}

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowPasswordHint(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	log.Printf("Entering password-hint with username: %s", username)

	if username == "" {
		http.Error(w, "No username provided, please type in your username first", http.StatusBadRequest)
		return
	}

	// Prepare the SQL query
	sqlQuery := "SELECT password_hint FROM users WHERE username = ?"
	log.Println(sqlQuery)
	var passwordHint string
	err := sqlite.DB.QueryRow(sqlQuery, username).Scan(&passwordHint)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No password found for "+username, http.StatusNotFound)
		} else {
			log.Println("Error querying database:", err)
			http.Error(w, "ERROR!", http.StatusInternalServerError)
		}
		return
	}

	if len(passwordHint) > 0 {
		formatString := fmt.Sprintf("Username '%s' has password: %s%s", username, passwordHint[:2], strings.Repeat("*", len(passwordHint)-2))
		log.Println(formatString)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`"%s"`, formatString)))
	} else {
		http.Error(w, "No password found for "+username, http.StatusNotFound)
	}

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

func updateInResponse(currentUser interface{}, w http.ResponseWriter) error {
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

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering ShowRegister")

	view.Render(w, "register.html", nil)
}
func ProcessRegister(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Entering ProcessRegister")
	username := r.FormValue("username")
	output := Output{Username: username}
	// This might be an error due to incorrect pointer logic
	if username == "" {
		output.Error = "No username provided, please type in your username first"
		view.Render(w, "register.html", output)
		return
	}
	log.Println("Creating database query")
	sqlQuery := "SELECT username FROM users WHERE username = '" + username + "'"
	log.Println(sqlQuery)
	row := sqlite.DB.QueryRow(sqlQuery)
	var err error
	var expectedUser string
	if err = row.Scan(expectedUser); err == sql.ErrNoRows {
		view.Render(w, "register-finish.html", output)
		return

	} else {
		//This is the case wehre row is not empty
		output.Error = "Username '" + username + "' already exists!"
		view.Render(w, "register.html", output)
		return
	}

}

func ShowRegisterFinish(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering ShowRegisterFinish")

	view.Render(w, "register-finish.html", nil)
}
func ProcessRegisterFinish(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	cpassword := r.FormValue("cpassword")
	realName := r.FormValue("realName")
	blabName := r.FormValue("blabName")
	output := Output{Username: username}
	if password != cpassword {
		log.Println("Password and Confirm Password do not match")
		output.Error = "The Password and Confirm Password values do not match. Please try again."
		view.Render(w, "register-finish.html", output)
		return
	}

	// // Execute the query

	//TODO: Test TOTP functionality
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "VeraDemo",
		AccountName: username,
	})
	if err != nil {
		log.Println("Failed to generate TOTP secret.")
		return
	}
	query := "insert into users (username, password, totp_secret, created_at, real_name, blab_name) values("
	query += "'" + username + "',"
	query += "'" + GetMD5Hash(password) + "',"
	query += "'" + secret.Secret() + "',"
	query += "datetime('now'),"
	query += "'" + realName + "',"
	query += "'" + blabName + "'"
	query += ");"
	log.Println(query)
	_, err = sqlite.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
	current_session := session.Instance(r)
	current_session.Values["username"] = username
	err = current_session.Save(r, w)
	if err != nil {
		log.Println("Couldn't set session value")
	}
	fmt.Println(current_session.Values["username"].(string))

	//view.Render(w, "feed.html", output)
	http.Redirect(w, r, "login?username="+username, http.StatusSeeOther)
}

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	// let type = req.query.type : Implement GET request reading
	output := Output{}
	log.Println("Entering ShowProfile")
	current_session := session.Instance(r)
	username := current_session.Values["username"].(string)

	// TODO: Fix page redirect without session
	// Currently an old session is getting picked up that should not be.
	fmt.Println(username)
	if username == "" {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=profile", http.StatusSeeOther)
		return
	}

	sqlMyHecklers := "SELECT users.username, users.blab_name, users.created_at FROM users LEFT JOIN listeners ON users.username = listeners.listener WHERE listeners.blabber=? AND listeners.status='Active';"
	log.Println(sqlMyHecklers)
	hecklers := []models.Blabber{}
	rows, err := sqlite.DB.Query(sqlMyHecklers, username)
	if err != nil {
		log.Println(err)
		view.Render(w, "profile.html", output)
		return
	}

	//Scans all results from query into the hecklers array
	for rows.Next() {
		i := models.Blabber{}
		err = rows.Scan(&i.Username, &i.BlabName, &i.CreatedDate)
		if err != nil {
			log.Println("Error scanning sql data response")
			view.Render(w, "profile.html", output)
			return
		}
		hecklers = append(hecklers, i)

	}
	events := []string{}
	sqlMyEvents := "select event from users_history where blabber=\"" + username + "\" ORDER BY eventid DESC; "
	log.Println(sqlMyEvents)
	rows, err = sqlite.DB.Query(sqlMyEvents)
	if err != nil {
		log.Println("Couldn't retreive event history")
		view.Render(w, "profile.html", output)
		return
	}

	for rows.Next() {
		var i string
		err = rows.Scan(&i)
		if err != nil {
			log.Println("Error scanning sql data response")
			view.Render(w, "profile.html", output)
			return
		}
		events = append(events, i)

	}

	sqlQuery := "SELECT username, real_name, blab_name, totp_secret FROM users WHERE username = '" + username + "'"
	log.Println(sqlQuery)

	row := sqlite.DB.QueryRow(sqlQuery)

	if err = row.Scan(&output.Username, &output.RealName, &output.BlabName, &output.TotpSecret); err == sql.ErrNoRows {
		output.Error = "Access Denied: no user data found"
		view.Render(w, "login.html", output)
		return

	}
	output.Events = events
	output.Hecklers = hecklers
	view.Render(w, "profile.html", output)

}

type JSONResponse struct {
	Message string
}

func ProcessProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering ProcessProfile")

	realName := r.FormValue("realName")
	blabName := r.FormValue("blabName")
	username := r.FormValue("username")
	//TODO: Check for supplied file

	current_session := session.Instance(r)
	sessionUsername := current_session.Values["username"].(string)
	if sessionUsername == "" {
		log.Println("User is not logged in - redirecting...")
		http.Redirect(w, r, "login?target=profile", http.StatusSeeOther)
		return
	}
	frame := JSONResponse{}
	//TODO: Print out user agent
	// log.Println("User is logged in - continuing... UA=" + )
	log.Println("user logged in")
	oldUsername := sessionUsername
	log.Println("Executing the update prepared statement")
	result, err := sqlite.DB.Exec("UPDATE users SET real_name=?, blab_name=? WHERE username=?;", realName, blabName, sessionUsername)
	// Nested statements for error handling
	if err != nil {
		log.Println("Error updating user")
	} else {
		RowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
		} else if RowsAffected != 1 {
			frame.Message = "<script>alert('An error occurred, please try again.');</script>"
			response, err := json.Marshal(frame)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}
	}

	//Rename profijle image if username changes
	if !(username == oldUsername) {
		exists := false
		newUsername := strings.ToLower(username)

		log.Println("Preparing the duplicate username check Prepared Statement")
		row := sqlite.DB.QueryRow("SELECT username FROM users WHERE username=?", newUsername)
		var err error
		if err = row.Scan(); err != sql.ErrNoRows {
			log.Println("Username: " + username + " already exists. Try again")
			exists = true

		}
		if exists {
			frame.Message = "<script>alert('That username already exists. Please try another.');</script>"
			response, err := json.Marshal(frame)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusConflict)
			os.Stdout.Write(response)
			w.Write(response)
			return
		}

		//attempt to update username
		oldUsername = strings.ToLower(oldUsername)
		log.Println("Creating Transaction")
		tx, err := sqlite.DB.BeginTx(context.Background(), nil)
		if err != nil {
			log.Println(err)
			return
		}
		//Call rollback when function is returned. If function returns early, transaction rolls back before being committed.
		defer tx.Rollback()

		sqlStrQueries := []string{
			"UPDATE users SET username='?' WHERE username='?'",
			"UPDATE blabs SET blabber='?' WHERE blabber='?'",
			"UPDATE comments SET blabber='?' WHERE blabber='?'",
			"UPDATE listeners SET blabber='?' WHERE blabber='?'",
			"UPDATE listeners SET listener='?' WHERE listener='?'",
			"UPDATE users_history SET blabber='?' WHERE blabber='?'"}

		log.Println("Executing Transactions")
		for i := 0; i < len(sqlStrQueries); i++ {
			_, err := tx.Exec(sqlStrQueries[i], newUsername, oldUsername)
			if err != nil {
				log.Println(err)
				return
			}
		}
		//Commit Transactions
		if err = tx.Commit(); err != nil {
			log.Println(err)
			return
		}
		//oldImage := GetProfileImageFromUsername(oldUsername)

	}

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// func GetProfileImageFromUsername(username){
// 	files :=
// }
