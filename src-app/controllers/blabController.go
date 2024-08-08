package controllers

import (
	"net/http"

	// "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	"verademo-go/src-app/shared/view"

	"github.com/ian-kent/go-log/log"
)

var sqlBlabsByMe = `SELECT blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs LEFT JOIN comments ON blabs.blabid = comments.blabid ` +
	`WHERE blabs.blabber = '%s' GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC;`

var sqlBlabsForMe = `SELECT users.username, users.blab_name, blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs INNER JOIN users ON blabs.blabber = users.username INNER JOIN listeners ON blabs.blabber = listeners.blabber ` +
	`LEFT JOIN comments ON blabs.blabid = comments.blabid WHERE listeners.listener = '%s' ` +
	`GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC LIMIT {} OFFSET {};`

func ShowFeed(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	username := sess.Values["username"].(string)

	if username == "" {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=feed", http.StatusFound)
		return
	}

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	log.Println("Executing query to get all 'Blabs for me'")
	//blabsForMe := fmt.Sprintf(sqlBlabsForMe, 10, 0)
	//db.Db.Exec(blabsForMe, username)
	// cursor.execute(blabsForMe % (username,))
	// blabsForMeResults = cursor.fetchall()

	view.Render(w, "feed.html", nil)

}
