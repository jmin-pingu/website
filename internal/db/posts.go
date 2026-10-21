package db

import (
	"context"
	"fmt"
	"internal/ds"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

const POST_UPDATE = `
	UPDATE posts
	SET tags=$1, title=$2, date=$3, display=$4, content=$5
	WHERE link = $6;
`

const POST_INSERT = `
	INSERT INTO posts (post_id, tags, title, link, date, display, content)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
`

type Post struct {
	PostID  string
	Tags    []string
	Title   string
	Link    string
	Display bool
	Date    time.Time
	Content string // Format: HTML
}

func InitPosts(dbpool *pgxpool.Pool) {
	var (
		stmt string
		err  error
		dat  []byte
	)

	// Read SQL schema for posts
	pwd, err := os.Getwd()
	dat, err = os.ReadFile(pwd + "/internal/db/posts_schema.sql")

	if err != nil {
		fmt.Fprintf(os.Stderr, "`InitPosts` failed to read schema: %v\n", err)
		os.Exit(1)
	}
	stmt = string(dat)

	// Execute script
	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`InitPosts` failed: %v\n", err)
		os.Exit(1)
	}
}

func GetPosts(dbpool *pgxpool.Pool) ([]*Post, error) {
	query := `SELECT * FROM posts WHERE display = 'true' ORDER BY date DESC;`

	var posts []*Post
	err := pgxscan.Select(context.Background(), dbpool, &posts, query)
	return posts, err
}

// GetPostByLink returns nil if no displayed post has that link
func GetPostByLink(dbpool *pgxpool.Pool, link string) (*Post, error) {
	query := `SELECT * FROM posts WHERE link = $1 AND display = 'true';`

	var post Post
	err := pgxscan.Get(context.Background(), dbpool, &post, query, link)
	if pgxscan.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func UploadPost(dbpool *pgxpool.Pool, cmd string, tags []string, title string, link string, date time.Time, display bool, content string) {
	var err error

	lower_tags := make([]string, len(tags))
	for i, v := range tags {
		lower_tags[i] = strings.ToLower(v)
	}

	// Values are passed as query parameters, so quotes in them can't break the SQL
	switch cmd {
	case "update":
		_, err = dbpool.Exec(context.Background(), POST_UPDATE,
			lower_tags, title, date, display, content, link)
	case "insert":
		_, err = dbpool.Exec(context.Background(), POST_INSERT,
			Hash(link).String(), lower_tags, title, link, date, display, content)
	default:
		panic("`UploadPost`: `cmd` should either be `update` or `insert`")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "`UploadPost` failed: %v\n", err)
		os.Exit(1)
	}
	if cmd == "update" {
		fmt.Printf("\tupdated: %v\n", link)
	} else {
		fmt.Printf("\tinserted: %v\n", link)
	}
}

// UploadToPosts returns the links it uploaded
func UploadToPosts(dbpool *pgxpool.Pool, paths []string) []string {
	links := make([]string, 0)
	// NOTE: check for appropriate thing to upload
	log.Println("--------- Uploading Paths to `posts` table ---------")
	log.Println("paths: %v", paths)
	// Parse through paths
	for _, path := range paths {
		// NOTE: a bunch of checks to make sure the file is valid
		if _, err := os.Stat(path); err == nil {
			fmt.Println("valid path:", path)
		} else if match, _ := regexp.MatchString(".md$", path); !match {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read arguments : file should be .md\n")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read sql template: %v\n", err)
			os.Exit(1)
		}

		// Read md file and extract content and metadata
		md, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read markdown file: %v", err)
		}
		content := mdToHTML(md)
		tags, title, link, display, date := parse_md_metadata(md)

		// Based on contents of md file, upload or insert
		if Exists(dbpool, "posts", "link", link) {
			UploadPost(dbpool, "update", tags, title, link, date, display, content)
		} else {
			UploadPost(dbpool, "insert", tags, title, link, date, display, content)
		}
		links = append(links, link)
	}
	return links
}

func parse_md_metadata(md []byte) ([]string, string, string, bool, time.Time) {
	// Get metadata from markdown
	metadata := getMetadata(md) // Consider what additional metadata I would want to consider

	// Parse date
	date, err := time.Parse("2006-01-02", metadata["date"].(string))
	if err != nil {
		log.Fatalf("failed to parse date from metadata: %v", err)
	}

	// Parse tags
	tags_set := make(ds.Set[string], 0)
	parsed_tags := strings.Split(metadata["tags"].(string), ",")
	for _, tag := range parsed_tags {
		t := strings.ToLower(strings.TrimSpace(tag))
		tags_set.Add(t) // Per post tag set: ds.Set
	}
	tags := make([]string, 0)
	for k, _ := range tags_set {
		tags = append(tags, k)
	}
	return tags, metadata["title"].(string), metadata["link"].(string), metadata["display"].(bool), date
}
