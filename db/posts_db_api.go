package db

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	Content string // Format: HTML
}

func InitPosts(dbpool *pgxpool.Pool, clean bool) {
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

	if clean {
		_, err = dbpool.Exec(context.Background(), `DROP TABLE posts;`)
	}
	// Execute script
	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`init_posts` failed: %v\n", err)
		os.Exit(1)
	}
}

// type Post struct {
// 	PostID  int
// 	Tags    []string
// 	Title   string
// 	Link    string
// 	Date    time.Time
// 	Content string // Format: HTML
// }

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

func UploadPost(dbpool *pgxpool.Pool, cmd string, tags []string, title string, link string, date time.Time, content string) {
	var (
		script string
		err    error
	)
	UPDATE_TEMPLATE := `
		UPDATE posts 
		SET tags=%s, title='%s', link='%s', date='%s', content=$html$%s$html$
		WHERE link = '%s';
	`

	INSERT_TEMPLATE := `
		INSERT INTO posts (tags, title, link, date, content)
		VALUES (
			%s,
			'%s',
			'%s',
			'%s',
			$html$%s$html$
		);
	`

	// Parse date and tags to make sure inputs work with SQL
	parsed_date := strings.Split(date.String(), " ")[0] // Change date to proper formatting for SQL

	parsed_tags := "ARRAY["
	for _, v := range tags {
		parsed_tags = parsed_tags + "'" + v + "',"
	}
	parsed_tags = parsed_tags[:len(parsed_tags)-1] + "]"

	switch cmd {
	case "update":
		script = fmt.Sprintf(UPDATE_TEMPLATE, parsed_tags, title, link, parsed_date, content, link)
		fmt.Printf("\tupdated link: %v\n", link)
	case "insert":
		script = fmt.Sprintf(INSERT_TEMPLATE, parsed_tags, title, link, parsed_date, content)
		fmt.Printf("\tupdated link: %v\n", link)
	default:
		panic("`UploadPost`: `cmd` should either be `update` or `insert`")
	}
	// Execute script
	_, err = dbpool.Exec(context.Background(), script)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`UploadPost` failed: %v\n", err)
		os.Exit(1)
	}
}
