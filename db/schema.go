package db

import "time"

type Post struct {
	PostID  int
	Tags    []string
	Title   string
	Link    string
	Date    time.Time
	Content string // Should be .md formated content
}
