package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"verademo-go/src-app/models"
	sqlite "verademo-go/src-app/shared/db"
	session "verademo-go/src-app/shared/session"
	"verademo-go/src-app/shared/view"

	"log"
)

var sqlBlabsByMe = `SELECT blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs LEFT JOIN comments ON blabs.blabid = comments.blabid ` +
	`WHERE blabs.blabber = ? GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC;`

var sqlBlabsForMe = `SELECT users.username, users.blab_name, blabs.content, blabs.timestamp, COUNT(comments.blabber), blabs.blabid ` +
	`FROM blabs INNER JOIN users ON blabs.blabber = users.username INNER JOIN listeners ON blabs.blabber = listeners.blabber ` +
	`LEFT JOIN comments ON blabs.blabid = comments.blabid WHERE listeners.listener = ? ` +
	`GROUP BY blabs.blabid ORDER BY blabs.timestamp DESC LIMIT %d OFFSET %d;`

func ShowFeed(w http.ResponseWriter, r *http.Request) {

	// Struct for variables to pass to the feed template
	type FeedVars struct {
		BlabsByOthers []models.Blab
		BlabsByMe     []models.Blab
		CurrentUser   string
		Error         string
	}

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=feed", http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	var outputs FeedVars

	// Set an error if one was given in response (usually taken from ProcessFeed)
	resError, err := r.Cookie("errorMsg")
	if err == nil {
		outputs.Error = resError.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "errorMsg",
			MaxAge: -1,
		})
	}

	// Get blabs from blabbers that are being listened to
	log.Println("Executing query to get all 'Blabs for me'")
	blabsForMe := fmt.Sprintf(sqlBlabsForMe, 10, 0)
	blabsForMeResults, err := sqlite.DB.Query(blabsForMe, username)
	if err != nil {
		errMsg := "Error getting 'Blabs for me': \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "feed.html", outputs)
		return
	}

	// Close the results object when they have been used up
	defer blabsForMeResults.Close()

	// Add each blab found to a variable to be passed to the template
	var feedBlabs []models.Blab

	for blabsForMeResults.Next() {
		var author models.Blabber
		var post models.Blab

		if err := blabsForMeResults.Scan(&author.Username, &author.BlabName, &post.Content, &post.PostDate, &post.CommentCount, &post.ID); err != nil {
			errMsg := "Error reading data from 'Blabs for me' query: \n" + err.Error()
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

	// Get blabs from the current user
	log.Println("Executing query to get all of user's Blabs")
	blabsByMeResults, err := sqlite.DB.Query(sqlBlabsByMe, username)
	if err != nil {
		errMsg := "Error getting 'Blabs for me': \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "feed.html", outputs)
		return
	}

	// Close the results object when they have been used up
	defer blabsByMeResults.Close()

	// Add each blab found to a variable to be passed to the template
	var myBlabs []models.Blab

	for blabsByMeResults.Next() {
		var post models.Blab

		if err := blabsByMeResults.Scan(&post.Content, &post.PostDate, &post.CommentCount, &post.ID); err != nil {
			errMsg := "Error reading data from 'Blabs by me' query: \n" + err.Error()
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
	countParam := r.URL.Query().Get("count")
	lenParam := r.URL.Query().Get("len")

	// Template for response
	template := "<li><div>" + "\t<div class=\"commenterImage\">" + "\t\t<img src=\"images/%s\">" +
		"\t</div>" + "\t<div class=\"commentText\">" + "\t\t<p>%s</p>" +
		"\t\t<span class=\"date sub-text\">by %s on %s</span><br>" +
		"\t\t<span class=\"date sub-text\"><a href=\"blab?blabid=%d\">%d Comments</a></span>" + "\t</div>" +
		"</div></li>"

	// Convert GET parameters to integers
	count, err := strconv.Atoi(countParam)
	if err != nil {
		log.Println("Error converting count:" + countParam + " to integer: \n" + err.Error())
		http.Redirect(w, r, "feed", http.StatusBadRequest)
		return
	}

	len, err := strconv.Atoi(lenParam)
	if err != nil {
		log.Println("Error converting len:" + lenParam + " to integer: \n" + err.Error())
		http.Redirect(w, r, "feed", http.StatusBadRequest)
		return
	}

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=feed", http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	// Run SQL query
	log.Println("Executing query to get more blabs")
	blabsForMe := fmt.Sprintf(sqlBlabsForMe, len, count)
	results, err := sqlite.DB.Query(blabsForMe, username)
	if err != nil {
		errMsg := "Error getting more blabs: \n" + err.Error()
		log.Println(errMsg)
		http.Redirect(w, r, "feed", http.StatusBadRequest)
		return
	}

	// Close the results object when they have been used up
	defer results.Close()

	// Add each blab found to the response using the template
	var ret string

	for results.Next() {
		var author models.Blabber
		var post models.Blab

		if err := results.Scan(&author.Username, &author.BlabName, &post.Content, &post.PostDate, &post.CommentCount, &post.ID); err != nil {
			errMsg := "Error reading data from 'more feed' query: \n" + err.Error()
			log.Println(errMsg)
			http.Redirect(w, r, "feed", http.StatusBadRequest)
			return
		}

		ret += fmt.Sprintf(template, author.GetProfileImageFromUsername(), post.Content, author.BlabName, models.Timestamp(post.PostDate), post.ID, post.CommentCount)

	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(ret))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}

}

func ProcessFeed(w http.ResponseWriter, r *http.Request) {
	blab := r.FormValue("blab")

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=feed", http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	// Post a blab
	log.Println("Executing query to post a blab")
	addBlabSql := "INSERT INTO blabs (blabber, content, timestamp) values (?, ?, datetime('now'));"
	result, err := sqlite.DB.Exec(addBlabSql, username, blab)
	if err != nil {
		errMsg := "Error posting blab: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "feed", http.StatusSeeOther)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		errMsg := "Error posting blab: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "feed", http.StatusSeeOther)
		return
	}
	if rows != 1 {
		errMsg := fmt.Sprintf("Expected to affect 1 row, affected %d.", rows)
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "feed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "feed", http.StatusSeeOther)
}

func ShowBlab(w http.ResponseWriter, r *http.Request) {

	// Struct for variables to pass to the blab template
	type BlabVars struct {
		Content  string
		BlabName string
		BlabID   int
		Comments []models.Comment
		Error    string
	}

	blabidParam := r.URL.Query().Get("blabid")

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=blab?blabid="+blabidParam, http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	blabDetailsSql := "SELECT blabs.content, users.blab_name " +
		"FROM blabs INNER JOIN users ON blabs.blabber = users.username " +
		"WHERE blabs.blabid = ?;"

	blabCommentsSql := "SELECT users.username, users.blab_name, comments.content, comments.timestamp " +
		"FROM comments INNER JOIN users ON comments.blabber = users.username " +
		"WHERE comments.blabid = ? ORDER BY comments.timestamp DESC;"

	var outputs BlabVars

	// Set an error if one was given in response (usually taken from ProcessBlab)
	resError, err := r.Cookie("errorMsg")
	if err == nil {
		outputs.Error = resError.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "errorMsg",
			MaxAge: -1,
		})
	}

	// Convert GET parameter to integer
	blabid, err := strconv.Atoi(blabidParam)
	if err != nil {
		errMsg := "Error converting blab ID " + blabidParam + " to integer. Check the blab ID: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "blab.html", outputs)
		return
	}

	outputs.BlabID = blabid

	// Get blabs from blabbers that are being listened to
	log.Println("Executing query to get blab details")
	blabDetailsResult := sqlite.DB.QueryRow(blabDetailsSql, blabid)
	err = blabDetailsResult.Scan(&outputs.Content, &outputs.BlabName)
	switch {
	case err == sql.ErrNoRows:
		errMsg := "No blab found with ID:" + blabidParam + " \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "blab.html", outputs)
		return
	case err != nil:
		errMsg := "Error getting blab details: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "blab.html", outputs)
		return
	}

	// Get the blab's comments
	log.Println("Executing query to get all comments for blab")
	blabCommentsResults, err := sqlite.DB.Query(blabCommentsSql, blabid)
	if err != nil {
		errMsg := "Error getting blab comments: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "blab.html", outputs)
		return
	}

	// Close the results object when they have been used up
	defer blabCommentsResults.Close()

	// Add each comment found to a variable to be passed to the template
	var comments []models.Comment

	for blabCommentsResults.Next() {
		var author models.Blabber
		var comment models.Comment

		if err := blabCommentsResults.Scan(&author.Username, &author.BlabName, &comment.Content, &comment.PostDate); err != nil {
			errMsg := "Error reading data from blab comments: \n" + err.Error()
			log.Println(errMsg)
			outputs.Error = errMsg
			view.Render(w, "blab.html", outputs)
			return
		}

		comment.Author = author
		comment.PostDate = models.Timestamp(comment.PostDate)

		comments = append(comments, comment)

	}

	outputs.Comments = comments

	view.Render(w, "blab.html", outputs)
}

func ProcessBlab(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	blabid := r.FormValue("blabid")

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=blab?blabid="+blabid, http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	// Post a comment
	log.Println("Executing query to post a comment")
	addCommentSql := "INSERT INTO comments (blabid, blabber, content, timestamp) values (?, ?, ?, datetime('now'));"
	result, err := sqlite.DB.Exec(addCommentSql, blabid, username, comment)
	if err != nil {
		errMsg := "Error posting comment: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "/blab?blabid="+blabid, http.StatusSeeOther)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		errMsg := "Error posting comment: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "/blab?blabid="+blabid, http.StatusSeeOther)
		return
	}
	if rows != 1 {
		errMsg := fmt.Sprintf("Expected to affect 1 row, affected %d.", rows)
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "/blab?blabid="+blabid, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blab?blabid="+blabid, http.StatusSeeOther)
}

func ShowBlabbers(w http.ResponseWriter, r *http.Request) {
	// Struct for variables to pass to the feed template
	type BlabbersVars struct {
		Blabbers []models.Blabber
		Error    string
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "blab_name ASC"
	}

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=blabbers", http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	var outputs BlabbersVars

	// Set an error if one was given in response (usually taken from ProcessBlabbers)
	resError, err := r.Cookie("errorMsg")
	if err == nil {
		outputs.Error = resError.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "errorMsg",
			MaxAge: -1,
		})
	}

	blabbersSql := "SELECT users.username," + " users.blab_name," + " users.created_at," +
		" SUM(iif(listeners.listener=?, 1, 0)) as listeners," +
		" SUM(iif(listeners.status='Active',1,0)) as listening" +
		" FROM users LEFT JOIN listeners ON users.username = listeners.blabber" +
		" WHERE users.username NOT IN ('admin','admin-totp',?)" + " GROUP BY users.username" + " ORDER BY " + sort + ";"

	// Get the list of blabbers
	log.Println("Executing query to get all blabbers")
	blabbersResults, err := sqlite.DB.Query(blabbersSql, username, username)
	if err != nil {
		errMsg := "Error getting blab comments: \n" + err.Error()
		log.Println(errMsg)
		outputs.Error = errMsg
		view.Render(w, "blabbers.html", outputs)
		return
	}

	// Close the results object when they have been used up
	defer blabbersResults.Close()

	// Add each blabber found to a variable to be passed to the template
	var blabbers []models.Blabber

	for blabbersResults.Next() {
		var blabber models.Blabber

		if err := blabbersResults.Scan(&blabber.Username, &blabber.BlabName, &blabber.CreatedDate, &blabber.NumberListeners, &blabber.NumberListening); err != nil {
			errMsg := "Error reading data from blabbers: \n" + err.Error()
			log.Println(errMsg)
			outputs.Error = errMsg
			view.Render(w, "blabbers.html", outputs)
			return
		}

		blabber.CreatedDate = models.Timestamp(blabber.CreatedDate)

		blabbers = append(blabbers, blabber)

	}

	outputs.Blabbers = blabbers

	view.Render(w, "blabbers.html", outputs)
}

func ProcessBlabbers(w http.ResponseWriter, r *http.Request) {
	blabberUsername := r.FormValue("blabberUsername")
	command := r.FormValue("command")

	// Check session username
	sess := session.Instance(r)

	if sess.Values["username"] == nil {
		log.Println("User is not Logged In - redirecting...")
		http.Redirect(w, r, "login?target=blabbers", http.StatusFound)
		return
	}

	username := sess.Values["username"].(string)

	log.Println("User is Logged In - continuing... UA=" + r.Header.Get("user-agent") + " U=" + username)

	if command == "" {
		errMsg := "Empty command provided..."
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusFound)
		return
	}

	sqlQuery := ""
	var eventAction string

	switch command {
	case "ignore":
		sqlQuery = "DELETE FROM listeners WHERE blabber=? AND listener=?;"
		eventAction = " is now ignoring "
	case "listen":
		sqlQuery = "INSERT INTO listeners (blabber, listener, status) values (?, ?, 'Active');"
		eventAction = " started listening to "
	default:
		errMsg := "Invalid Command"
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusFound)
		return
	}

	// Execute the command
	log.Println("Executing command")
	result, err := sqlite.DB.Exec(sqlQuery, blabberUsername, username)
	if err != nil {
		errMsg := "Error executing command: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		errMsg := "Error executing command: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}
	if rows != 1 {
		errMsg := fmt.Sprintf("Expected to affect 1 row, affected %d.", rows)
		log.Println(errMsg)
		w.Header().Add("errorMsg", errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}

	// Get blab name of target user
	var blabName string
	sqlQuery = "SELECT blab_name FROM users WHERE username = '" + blabberUsername + "'"
	log.Println("Executing query to get blab name of the target user")
	queryResult := sqlite.DB.QueryRow(sqlQuery)
	err = queryResult.Scan(&blabName)
	switch {
	case err == sql.ErrNoRows:
		errMsg := "No user found with username:" + blabberUsername + " \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	case err != nil:
		errMsg := "Error getting target's blab name: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}

	// Insert into event history
	event := username + eventAction + blabberUsername + " (" + blabName + ")"
	sqlQuery = "INSERT INTO users_history (blabber, event) VALUES ('" + username + "', '" + event + "')"
	result, err = sqlite.DB.Exec(sqlQuery)
	if err != nil {
		errMsg := "Error adding event into history: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}
	rows, err = result.RowsAffected()
	if err != nil {
		errMsg := "Error adding event into history: \n" + err.Error()
		log.Println(errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}
	if rows != 1 {
		errMsg := fmt.Sprintf("Expected to affect 1 row, affected %d.", rows)
		log.Println(errMsg)
		w.Header().Add("errorMsg", errMsg)
		http.SetCookie(w, &http.Cookie{
			Name:  "errorMsg",
			Value: errMsg,
		})
		http.Redirect(w, r, "blabbers", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "blabbers", http.StatusSeeOther)
}
