package db

import (
	"bytes"
	"context"
	"fmt"
	"internal/ds"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Post struct {
	PostID  int
	Tags    []string
	Title   string
	Link    string
	Display bool
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
	dat, err = os.ReadFile(pwd + "/internal/db/posts_schema.sql")

	if err != nil {
		fmt.Fprintf(os.Stderr, "`InitPosts` failed to read schema: %v\n", err)
		os.Exit(1)
	}
	stmt = string(dat)

	if clean {
		_, err = dbpool.Exec(context.Background(), `DROP TABLE posts;`)
	}
	// Execute script
	_, err = dbpool.Exec(context.Background(), stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "`InitPosts` failed: %v\n", err)
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

func UploadPost(dbpool *pgxpool.Pool, cmd string, tags []string, title string, link string, date time.Time, display bool, content string) {
	var (
		script string
		err    error
	)
	UPDATE_TEMPLATE := `
		UPDATE posts 
		SET tags=%s, title='%s', link='%s', date='%s', display='%s', content=$html$%s$html$
		WHERE link = '%s';
	`

	INSERT_TEMPLATE := `
		INSERT INTO posts (tags, title, link, date, display, content)
		VALUES (
			%s,
			'%s',
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
		parsed_tags = parsed_tags + "'" + strings.ToLower(v) + "',"
	}
	parsed_tags = parsed_tags[:len(parsed_tags)-1] + "]"

	switch cmd {
	case "update":
		script = fmt.Sprintf(UPDATE_TEMPLATE, parsed_tags, title, link, parsed_date,
			strconv.FormatBool(display), content, link)

		fmt.Printf("\tupdated: %v\n", link)
	case "insert":
		script = fmt.Sprintf(INSERT_TEMPLATE, parsed_tags, title, link, parsed_date,
			strconv.FormatBool(display), content)
		fmt.Printf("\tinserted: %v\n", link)
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

func UploadToPosts(dbpool *pgxpool.Pool, paths []string) {
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
	}
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

// ----- Markdown-to-HTML Functions  -----
func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func mdToHTML(md []byte) string {
	var buf bytes.Buffer
	custom_parser := goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			meta.Meta,
			extension.Strikethrough,
			extension.Table,
			mathjax.MathJax,
		),
	)
	if err := custom_parser.Convert([]byte(md), &buf); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	return buf.String()
}

func getMetadata(md []byte) map[string]interface{} {
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	context := parser.NewContext()
	if err := markdown.Convert([]byte(md), &buf, parser.WithContext(context)); err != nil {
		log.Fatalf("failed to extract metadata: %v", err)
	}
	metadata := meta.Get(context)
	return metadata
}
