package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"mywebsite/components/pages"
	"mywebsite/ds"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func SetupRenders(e *echo.Echo) {
	// Render static pages
	e.GET("/", homeRenderer(&PAGES_METADATA))
	e.GET("/blog/", blogRenderer(&PAGES_METADATA, &POSTS_METADATA, &POSTS_TAGS, DISPLAY_TAGS))
	e.GET("/resources/", resourcesRenderer(&PAGES_METADATA))
	e.GET("/projects/", projectsRenderer(&PAGES_METADATA))
	e.GET("/creative/", projectsRenderer(&PAGES_METADATA))
	// Render blog posts
	RenderBlogPosts(e)
}

// NOTE: need to rethink how this will work with a db.
func RenderBlogPosts(e *echo.Echo) {
	dir, _ := os.Getwd()
	post_fnames, _ := os.ReadDir(dir + "/posts")
	for _, fname := range post_fnames {
		fmt.Printf("%v\n", fname.Name())
		path := "posts/" + fname.Name()
		md, err := os.ReadFile(path)
		if err != nil {
			log.Printf("failed to read markdown file: %v", err)
			continue
		}
		if strings.HasPrefix(fname.Name(), ".") {
			continue
		}
		// Get metadata from markdown
		metadata := getMetadata(md) // Consider what additional metadata I would want to consider

		date, err := time.Parse("2006-01-02", metadata["date"].(string))
		if err != nil {
			log.Fatalf("failed to parse date from metadata: %v", err)
		}

		tags := make(ds.Set[string], 0)
		parsed_tags := strings.Split(metadata["tags"].(string), ",")
		for _, tag := range parsed_tags {
			t := strings.ToLower(strings.TrimSpace(tag))
			POSTS_TAGS.Add(t) // Universal tag set: ds.OrderedList
			tags.Add(t)       // Per post tag set: ds.Set
		}

		// If there is a link argument take it, if not use a random uuid
		link := metadata["link"]
		var url string
		if link == nil {
			url = "/blog/" + uuid.NewString()
		} else {
			url = "/blog/" + link.(string)
		}
		POSTS_METADATA.AddPost(metadata["title"].(string), date, url, tags)

		// Convert markdown to HTML
		html := mdToHTML(md)
		e.GET(url, blogPageRenderer(metadata["title"].(string), html, &PAGES_METADATA, &POSTS_METADATA))
	}
}

// ----- Renderers -----
func homeRenderer(pages_metadata *ds.Pages) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.HomePage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func blogRenderer(pages_metadata *ds.Pages, posts_metadata *ds.Posts, tags *ds.OrderedList[string], filter_tags ds.Set[string]) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.BlogPage(pages_metadata, posts_metadata, tags, filter_tags).Render(context.Background(), c.Response().Writer)
	}
}

func resourcesRenderer(pages_metadata *ds.Pages) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.ResourcesPage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func projectsRenderer(pages_metadata *ds.Pages) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.ProjectsPage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func blogPageRenderer(fname, url string, pages_metadata *ds.Pages, posts_metadata *ds.Posts) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.BlogEntryPage(fname, url, pages_metadata, posts_metadata).Render(context.Background(), c.Response().Writer)
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
