package controllers

import (
	"fmt"
	"net/http"
	session "verademo-go/src-app/shared/session"
	"verademo-go/src-app/shared/view"
)

func ShowFeed(w http.ResponseWriter, r *http.Request) {
	current_session := session.Instance(r)

	fmt.Println(current_session.Values)
	view.Render(w, "feed.html", nil)
}
