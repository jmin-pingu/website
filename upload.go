package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mywebsite/db"
	"os"
	"regexp"

	"github.com/a-h/templ"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func main() {
	paths := os.Args[1:]
	// TODO: add functionality with glob and paths

	// Initialize db
	dbpool, err := db.GetConnection("websitedb")
	defer dbpool.Close() // Make sure to finish the transaction

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to websitedb: %v\n", err)
		os.Exit(1)
	}

	// Parse through paths
	fmt.Println("Input paths: %s", paths)
	for _, path := range paths {
		// NOTE: a bunch of checks to make sure the file is valid
		if _, err := os.Stat(path); err == nil {
			fmt.Println("Path is valid:", path)
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

		// Read md file
		md, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read markdown file: %v", err)
		}
		// html := mdToHTML(md)
		metadata := getMetadata(md)
		fmt.Println("Path %s Metadata %s", path, metadata)

		// Based on contents of md file, upload or insert
		// if db.Exists(dbpool, "posts", "link", ...) {
		// 	// TODO: call INSERT
		// } else {
		// 	// TODO: call UPDATE
		// }

	}
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
