package controllers

import (
	"net/http"
	"verademo-go/src-app/shared/view"
)

var sqlBlabsByMe = `SELECT blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs LEFT JOIN comments ON blabs.blabid = comments.blabid ` +
	`WHERE blabs.blabber = '%s' GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC;`

var sqlBlabsForMe = `SELECT users.username, users.blab_name, blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs INNER JOIN users ON blabs.blabber = users.username INNER JOIN listeners ON blabs.blabber = listeners.blabber ` +
	`LEFT JOIN comments ON blabs.blabid = comments.blabid WHERE listeners.listener = '%s' ` +
	`GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC LIMIT {} OFFSET {};`

func ShowFeed(w http.ResponseWriter, r *http.Request) {

	view.Render(w, "feed.html", nil)
}
