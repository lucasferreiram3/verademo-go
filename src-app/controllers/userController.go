package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"verademo-go/src-app/models"
	sqlite "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	"verademo-go/src-app/shared/utils"
	view "verademo-go/src-app/shared/view"

	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/pquerna/otp/totp"
)

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

	type LoginVars struct {
		Target string
		Error  string
	}

	target := req.URL.Query().Get("target")
	username := req.URL.Query().Get("username")

	var outputs LoginVars

	outputs.Target = target

	current_session := session.Instance(req)
	if current_session.Values["username"] != nil {
		log.Println("User is already logged in - redirecting...")
		if target == "" || target == "login" {
			target = "/feed"
		}
		http.Redirect(w, req, target, http.StatusFound)
		return
	}

	// Set an error if one was given in response (usually taken from ProcessLogin)
	resError, err := req.Cookie("errorMsg")
	if err == nil {
		outputs.Error = resError.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "errorMsg",
			MaxAge: -1,
		})
	}

	user, err := createFromRequest(req)
	if err != nil {
		errMsg := "Error creating user from request:" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "login.html", outputs)
		return
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else if user != nil {
		current_session.Values["username"] = user.Username
		_ = current_session.Save(req, w)
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

	log.Println("Entering showLogin with username " + username + " and target " + target)

	view.Render(w, "login.html", outputs)
}

func ProcessLogin(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering processLogin")

	username := req.FormValue("username")
	password := req.FormValue("password")
	remember := req.FormValue("remember")
	target := req.FormValue("target")

	var nextView string

	if target == "" {
		nextView = "feed"
	} else {
		nextView = target
	}

	// Check inputs before processing Query
	log.Println("Username: " + username + " Password: " + password)
	// Constructing SQL Query, using COALESECE in case of null password hints (new totp users that are manually registered were running into sql errors, this rectifies that)
	sqlQuery := "SELECT username, COALESCE(password_hint, '') as password_hint,  created_at, last_login, real_name, blab_name FROM users WHERE username = ? AND password = ?"

	var result models.User

	err := sqlite.DB.QueryRow(sqlQuery, username, utils.GetMD5Hash(password)).Scan(
		&result.Username,
		&result.PasswordHint,
		&result.CreatedDate,
		&result.LastLogin,
		&result.RealName,
		&result.BlabName,
	)
	log.Println("After Query: " + result.Username)
	// In case user does not exist
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "Username or password incorrect. Please try again."
			log.Println(errMsg)
			http.SetCookie(w, &http.Cookie{
				Name:  "errorMsg",
				Value: errMsg,
			})
			http.Redirect(w, req, "/login?target="+target, http.StatusSeeOther)
			return
		}
		errMsg := "An error has occurred. Please try again. \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, req, "/login?target="+target, http.StatusSeeOther)
		return
		// http.Error(w, "An error has occured", http.StatusInternalServerError)
	}

	log.Println("User found.")

	// Updating last login time
	_, err = sqlite.DB.Exec("UPDATE users SET last_login=datetime('now') WHERE username='" + result.Username + "';")
	if err != nil {
		errMsg := "Error updating last login for user: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, req, "/login?target="+target, http.StatusSeeOther)
		return
		// http.Error(w, "An error occurred", http.StatusInternalServerError)
	}

	// TOTP Handling
	if len(username) >= 4 && username[len(username)-4:] == "totp" {
		log.Println("User " + username + " Has TOTP Enabled!")
		current_session, _ := session.Store.Get(req, session.Name)
		current_session.Values["totp_username"] = result.Username
		current_session.Save(req, w)
		nextView = "/totp?totp_username=" + result.Username
	} else {
		// Handling the "remember me"
		if remember != "" {
			updateInResponse(result, w)
		}

		// Set session username
		current_session := session.Instance(req)
		current_session.Values["username"] = result.Username
		err = current_session.Save(req, w)
		if err != nil {
			errMsg := "Failed to set session value: \n" + err.Error()
			log.Println(errMsg)
			http.SetCookie(w, &http.Cookie{
				Name:  "errorMsg",
				Value: errMsg,
			})
			http.Redirect(w, req, "/login?target="+target, http.StatusSeeOther)
			return
		}
		log.Println("Setting session username to: " + username)
	}

	log.Println("Redirecting to view: " + nextView)
	http.Redirect(w, req, nextView, http.StatusSeeOther)

}
func ShowTotp(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("totp_username")

	if username == "" {
		log.Println("Username not found, Redirecting to login")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	var totpSecret string
	log.Println("Entering ShowTotp for : " + username)

	sqlQuery := "SELECT totp_secret FROM users WHERE username = ?"

	err := sqlite.DB.QueryRow(sqlQuery, username).Scan(&totpSecret)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User (TOTP) not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "An error has occurred", http.StatusInternalServerError)
		return
	}

	// If no secret generate TOTP
	if totpSecret == "" {
		secret, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "VeraDemo",
			AccountName: username,
		})
		if err != nil {
			log.Println("Error updating TOTP secret for user: ", err)
			http.Error(w, "An error occurred", http.StatusInternalServerError)
			return
		}
		totpSecret = secret.Secret()
		_, err = sqlite.DB.Exec("UPDATE users SET totp_secret = ? WHERE username = ?", totpSecret, username)
		if err != nil {
			log.Println("Error updating TOTP secret for user: ", err)
			http.Error(w, "An error occurred", http.StatusInternalServerError)
			return
		}
	}

	// Generate TOTP code, pass it to form
	data := map[string]interface{}{
		"TotpSecret": totpSecret,
		"Username":   username,
	}

	view.Render(w, "totp.html", data)
}

func ProcessTotp(w http.ResponseWriter, req *http.Request) {
	totpCode := req.FormValue("totpCode")
	current_session, _ := session.Store.Get(req, session.Name)
	username, ok := current_session.Values["totp_username"].(string)
	log.Println("Entering ProcessTotp with username: " + username + " and code: " + totpCode)

	if !ok || username == "" {
		log.Println("Username not found (TOTP), Redirecting to login...")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	var nextView string
	var totpSecret string
	nextView = "/login"

	sqlQuery := "SELECT totp_secret FROM users WHERE username = ?"

	err := sqlite.DB.QueryRow(sqlQuery, username).Scan(&totpSecret)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Failed to find TOTP secret in the database")
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		log.Println("An Database Error has occured: ", err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
	if totp.Validate(totpCode, totpSecret) {
		log.Println("TOTP validation success")
		session, _ := session.Store.Get(req, session.Name)
		session.Values["username"] = username
		session.Save(req, w)
		nextView = "/feed"
	} else {
		log.Println("TOTP validation failure!")
		session, _ := session.Store.Get(req, session.Name)
		session.Values["username"] = nil
		session.Values["totp_username"] = nil
		session.Save(req, w)
		nextView = "/login"
	}
	log.Println("Redirecting to view (TOTP): " + nextView)

	http.Redirect(w, req, nextView, http.StatusSeeOther)

}

func ProcessLogout(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering ProcessLogout")

	current_session, _ := session.Store.Get(req, session.Name)

	// Clearing session values
	current_session.Options.MaxAge = -1
	current_session.Save(req, w)

	// Delete cookies
	http.SetCookie(w, &http.Cookie{Name: "username", MaxAge: -1, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: session.Name, MaxAge: -1, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "user", MaxAge: -1, Path: "/"})

	log.Println("Successfully logged out, redirecting to login page....")
	http.Redirect(w, req, "/login", http.StatusSeeOther)

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
	query += "'" + utils.GetMD5Hash(password) + "',"
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
	output.Image = GetProfileImageFromUsername(output.Username)
	output.Events = events
	output.Hecklers = hecklers
	view.Render(w, "profile.html", output)

}

type JSONResponse struct {
	Message string
}

func ProcessProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering ProcessProfile")
	// Read in form values
	realName := r.FormValue("realName")
	blabName := r.FormValue("blabName")
	username := r.FormValue("username")
	// in, header, fileErr := req.FormFile("file")

	// defer in.Close()
	// out, err := os.OpenFile(header.Filename)

	//set directory for images
	dir := "/resources/images/"

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
				frame.Message = "<script>alert('Database transactions failed');</script>"
				response, e := json.Marshal(frame)
				if e != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusConflict)
				os.Stdout.Write(response)
				w.Write(response)
				return
			}
		}
		//Commit Transactions
		if err = tx.Commit(); err != nil {
			log.Println(err)
			return
		}
		oldImage := GetProfileImageFromUsername(oldUsername)
		if oldImage != "" {

			extension := oldImage[strings.LastIndex(oldImage, "."):]
			newImage := newUsername + extension
			log.Println("Renaming profile image from " + oldImage + " to " + newImage)

			e := os.Rename(filepath.Join(dir, oldImage), filepath.Join(dir, newImage))
			if e != nil {
				frame.Message = "<script>alert('An error occurred, please try again.');</script>"
				response, err := json.Marshal(frame)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusConflict)
				os.Stdout.Write(response)
				w.Write(response)
				log.Fatal(e)
			}

		}
		//Update session and cookie logic
		current_session.Values["username"] = username
		//TODO: Update cookie logic and remember me func.

	}
	// Update user profile image
	// if (fileErr == nil) {
	// 	oldImage := GetProfileImageFromUsername(username)
	// 	if oldImage != "" {
	// 		err := os.Remove(filepath.Join(dir,oldImage))
	// 		if err != nil {
	// 			log.Println("failed to remove existing file")
	// 		}
	// 	}
	// }

}

func GetProfileImageFromUsername(username string) string {
	imageFiles := os.DirFS("/resources/images")
	_, err := fs.ReadFile(imageFiles, username+".png")
	if err != nil {
		return ""
	}
	return username + ".png"
}

func createFromRequest(req *http.Request) (*models.User, error) {
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

	var user models.User
	if err := json.Unmarshal(decoded, &user); err != nil {
		log.Println("Error unmarshaling user from cookie:", err)
		return nil, err
	}

	log.Println("Username is:", user.Username)
	return &user, nil
}

func updateInResponse(currentUser models.User, w http.ResponseWriter) error {
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
