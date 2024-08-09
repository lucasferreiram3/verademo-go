package controllers

import (
	"fmt"
	"net/http"
	"verademo-go/src-app/models"
	sqlite "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	"verademo-go/src-app/shared/view"

	"github.com/ian-kent/go-log/log"
)

var sqlBlabsByMe = `SELECT blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs LEFT JOIN comments ON blabs.blabid = comments.blabid ` +
	`WHERE blabs.blabber = ? GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC;`

var sqlBlabsForMe = `SELECT users.username, users.blab_name, blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs INNER JOIN users ON blabs.blabber = users.username INNER JOIN listeners ON blabs.blabber = listeners.blabber ` +
	`LEFT JOIN comments ON blabs.blabid = comments.blabid WHERE listeners.listener = ? ` +
	`GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC LIMIT %d OFFSET %d;`

func ShowFeed(w http.ResponseWriter, r *http.Request) {

	type Outputs struct {
		BlabsByOthers []models.Blab
		BlabsByMe     []models.Blab
		CurrentUser   string
		Error         string
	}

	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=feed", http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	var outputs Outputs

	log.Println("Executing query to get all 'Blabs for me'")
	blabsForMe := fmt.Sprintf(sqlBlabsForMe, 10, 0)
	blabsForMeResults, err := sqlite.DB.Query(blabsForMe, username)
	if err != nil {
		errMsg := "Error getting 'Blabs for me':\n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "feed.html", outputs)
		return
	}

	defer blabsForMeResults.Close()

	var feedBlabs []models.Blab

	for blabsForMeResults.Next() {
		var author models.Blabber
		var post models.Blab

		if err := blabsForMeResults.Scan(&author.Username, &author.BlabName, &post.Content, &post.PostDate, &post.CommentCount, &post.Id); err != nil {
			errMsg := "Error reading data from 'Blabs for me' query:\n" + err.Error()
			log.Println(errMsg)
			outputs.Error = errMsg
			view.Render(w, "feed.html", outputs)
			return
		}

		post.Author = author
		post.PostDate = models.Timestamp(post.PostDate)

		feedBlabs = append(feedBlabs, post)

	}

	outputs.BlabsByOthers = feedBlabs
	outputs.CurrentUser = username

	log.Println("Executing query to get all of user's Blabs")
	blabsByMeResults, err := sqlite.DB.Query(sqlBlabsByMe, username)
	if err != nil {
		errMsg := "Error getting 'Blabs for me':\n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "feed.html", outputs)
		return
	}

	defer blabsByMeResults.Close()

	var myBlabs []models.Blab

	for blabsByMeResults.Next() {
		var post models.Blab

		if err := blabsByMeResults.Scan(&post.Content, &post.PostDate, &post.CommentCount, &post.Id); err != nil {
			errMsg := "Error reading data from 'Blabs by me' query:\n" + err.Error()
			log.Println(errMsg)
			outputs.Error = errMsg
			view.Render(w, "feed.html", outputs)
			return
		}

		post.PostDate = models.Timestamp(post.PostDate)

		myBlabs = append(myBlabs, post)

	}

	outputs.BlabsByMe = myBlabs

	view.Render(w, "feed.html", outputs)

}

func MoreFeed(w http.ResponseWriter, r *http.Request) {

}
