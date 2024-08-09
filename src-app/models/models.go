package models

import "time"

func Timestamp(postDate string) string {
	date, err := time.Parse(time.RFC3339, postDate)
	timestamp := date.Format("Jan 2, 2006")
	if err != nil {
		timestamp = "Error getting date:\n" + err.Error()
	}
	return timestamp
}
