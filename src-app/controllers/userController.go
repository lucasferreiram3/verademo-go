package controllers

import (
	"fmt"
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
}

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering ShowRegister")
	p := &Account{Error: "Bad"}
	view.Render(w, "register.html", p)
}
func ProcessRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	output := Output{username: username}
	fmt.Println("Entering ProcessRegister")
	fmt.Println("Creating session")
	current_session := session.Instance(r)
	current_session.Values["username"] = r.FormValue("username")
	fmt.Println(current_session.Values)

	view.Render(w, "register.html", &output)
}
