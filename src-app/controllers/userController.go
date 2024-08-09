package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	sqlite "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	view "verademo-go/src-app/shared/view"

	"github.com/pquerna/otp/totp"
)

type Account struct {
	Error    string
	Username string
}
type Output struct {
	Username string
	Error    string
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
	// fmt.Println(row)
	// fmt.Println("Creating session")
	// current_session := session.Instance(r)
	// current_session.Values["username"] = r.FormValue("username")
	// err = current_session.Save(r, w)
	// if err != nil {
	// 	log.Println("session error")
	// }
	// fmt.Println(current_session.Values["username"].(string))

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

	view.Render(w, "feed.html", output)
	return
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
