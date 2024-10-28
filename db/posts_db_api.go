package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	PostID  int
	Tags    []string
	Title   string
	Link    string
	Date    time.Time
	Content string // Should be .md formated content
}

func InitPosts(dbpool *pgxpool.Pool) {
	var (
		stmt string
		err  error
		dat  []byte
	)
	// Read SQL schema for posts

	pwd, err := os.Getwd()
	dat, err = os.ReadFile(pwd + "/db/posts_schema.sql")

	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_posts` failed to read schema: %v\n", err)
		os.Exit(1)
	}
	stmt = string(dat)

	_, err = dbpool.Exec(context.Background(), `DROP TABLE posts;`)
	// Execute script
	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_posts` failed: %v\n", err)
		os.Exit(1)
	}
}

func GetPosts(dbpool *pgxpool.Pool) []*Post {
	query := `SELECT * FROM posts;`

	var posts []*Post
	err := pgxscan.Select(context.Background(), dbpool, &posts, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`GetPosts` failed: %v\n", err)
		os.Exit(1)
	}
	return posts
}

const UPLOAD_TEMPLATE = `
	INSERT INTO posts (tags, title, link, date, content)
	VALUES (
		%s,
		%s,
		%s,
		%s,
		%s,
	);
`

const UPDATE_TEMPLATE = `
	UPDATE posts 
	SET tags=%s, title=%s, link=%s, date=%s, content=%s
	WHERE link == %s;
`
