package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mywebsite/db"
	"mywebsite/ds"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/a-h/templ"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func main() {
	paths := os.Args[1:]
	// Sanity check
	if len(paths) == 0 {
		panic("Did not provide paths to upload (.md files)")
	}

	// Initialize db
	dbpool, err := db.GetConnection("websitedb")
	defer dbpool.Close() // Make sure to finish the transaction

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to websitedb: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--------- Uploading Paths to `posts` table ---------")
	fmt.Println("Paths: %v", paths)
	// Parse through paths
	for _, path := range paths {
		// NOTE: a bunch of checks to make sure the file is valid
		if _, err := os.Stat(path); err == nil {
			fmt.Println("valid path:", path)
		} else if match, _ := regexp.MatchString(".md$", path); !match {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read arguments : file should be .md\n")
			os.Exit(1)
		} else if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "`upload.go` failed to read sql template: %v\n", err)
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
		tags, title, link, date := parse_md_metadata(md)

		// Based on contents of md file, upload or insert
		if db.Exists(dbpool, "posts", "link", link) {
			db.UploadPost(dbpool, "update", tags, title, link, date, content)
		} else {
			db.UploadPost(dbpool, "insert", tags, title, link, date, content)
		}
	}
	fmt.Println("success 🎉")
}

func parse_md_metadata(md []byte) ([]string, string, string, time.Time) {
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
	return tags, metadata["title"].(string), metadata["link"].(string), date
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
