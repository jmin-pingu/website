package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func mdToHTML(md []byte) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
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
		rlog.Fatalf("failed to extract metadata: %v", err)
	}
	metadata := meta.Get(context)
	fmt.Printf("Type of metadata %T\n", metadata)
	return metadata
}

func BlogEntryPage(fname, url string) {
	md, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalf("failed to read markdown file: %v", err)
	}
	html := mdToHTML(md)
	metadata := getMetadata(md)
	for key, val := range metadata {
		// Add appropriate metadata to POSTS_METADATA
	}
}
