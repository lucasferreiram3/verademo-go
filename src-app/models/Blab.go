package models

type Blab struct {
	Id           int    `db:"blabid"`
	Content      string `db:"content"`
	PostDate     string `'db:"timestamp"`
	CommentCount int
	Author       Blabber `'db:"blabber"`
}
