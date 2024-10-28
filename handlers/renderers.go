package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"mywebsite/components/pages"
	"mywebsite/db"
	"mywebsite/ds"

	"github.com/a-h/templ"
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
	dbpool, err := db.GetConnection("websitedb")
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	posts := db.GetPosts(dbpool)

	for _, post := range posts {
		tags := make(ds.Set[string], 0)
		for _, tag := range post.Tags {
			t := strings.ToLower(strings.TrimSpace(tag))
			POSTS_TAGS.Add(t) // Universal tag set: ds.OrderedList
			tags.Add(t)       // Per post tag set: ds.Set
		}

		// If there is a link argument take it, if not use a random uuid
		url := "/blog/" + post.Link
		POSTS_METADATA.AddPostMetadata(post.Title, post.Date, post.PostID, url, tags)

		// Convert markdown to HTML
		html := mdToHTML([]byte(post.Content))
		e.GET(url, blogPageRenderer(post.Title, html, &PAGES_METADATA, &POSTS_METADATA))
	}
}

// ----- Renderers -----
func homeRenderer(pages_metadata *ds.PagesMetadata) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.HomePage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func blogRenderer(pages_metadata *ds.PagesMetadata, posts_metadata *ds.PostsMetadata, tags *ds.OrderedList[string], filter_tags ds.Set[string]) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.BlogPage(pages_metadata, posts_metadata, tags, filter_tags).Render(context.Background(), c.Response().Writer)
	}
}

func resourcesRenderer(pages_metadata *ds.PagesMetadata) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.ResourcesPage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func projectsRenderer(pages_metadata *ds.PagesMetadata) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.ProjectsPage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func blogPageRenderer(fname, url string, pages_metadata *ds.PagesMetadata, posts_metadata *ds.PostsMetadata) echo.HandlerFunc {
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
