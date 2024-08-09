package controllers

import (
	"net/http"
	view "verademo-go/src-app/shared/view"
)

func ShowReset(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "reset.html", nil)
}

func ProcessReset(w http.ResponseWriter, r *http.Request) {

}
