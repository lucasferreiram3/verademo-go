package controllers

import (
	"net/http"
	"verademo-go/src-app/shared/view"
)

func ShowFeed(w http.ResponseWriter, r *http.Request) {

	view.Render(w, "feed.html", nil)
}
