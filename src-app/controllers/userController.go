package controllers

import (
	"fmt"
	"log"
	"net/http"
	session "verademo-go/src-app/shared/session"
	view "verademo-go/src-app/shared/view"
)

type Account struct {
	Error    string
	Username string
}
type Output struct {
	username string
	Error    string
}

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering ShowRegister")
	p := &Account{Error: "Bad"}
	view.Render(w, "register.html", p)
}
func ProcessRegister(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Entering ProcessRegister")
	username := r.FormValue("username")
	output := Output{username: username}
	// This might be an error due to incorrect pointer logic
	if username != "" {
		output.Error = "No username provided, please type in your username first"
		view.Render(w, "register.html", output)
		return
	}
	fmt.Println()

	fmt.Println("Creating session")
	current_session := session.Instance(r)
	current_session.Values["username"] = r.FormValue("username")
	err := current_session.Save(r, w)
	if err != nil {
		log.Println("session error")
	}
	fmt.Println(current_session.Values["username"].(string))

	view.Render(w, "register.html", &output)
}
