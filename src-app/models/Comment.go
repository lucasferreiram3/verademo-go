package models

type Comment struct {
	CommentID int
	BlabID    int
	Author    Blabber
	Content   string
	PostDate  string
}
