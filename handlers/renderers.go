package handlers

import (
	"context"
	"log"
	"os"
	"strings"

	"mywebsite/components/pages"
	"mywebsite/db"
	"mywebsite/ds"

	"github.com/labstack/echo/v4"
)

func SetupRenders(e *echo.Echo) {
	// Render static pages
	e.GET("/", homeRenderer(&PAGES_METADATA))
	e.GET("/blog/", blogRenderer(&PAGES_METADATA, &POSTS_METADATA, &POSTS_TAGS, DISPLAY_TAGS))
	e.GET("/resources/", resourcesRenderer(&PAGES_METADATA))
	e.GET("/projects/", projectsRenderer(&PAGES_METADATA))
	e.GET("/creative/", creativeRenderer(&PAGES_METADATA))
	e.GET("/reading_list/", readingListRenderer(&PAGES_METADATA, &BOOKS))
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
		// TODO: fix functionality
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
		e.GET(url, blogPageRenderer(post.Title, post.Content, &PAGES_METADATA, &POSTS_METADATA))
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

func creativeRenderer(pages_metadata *ds.PagesMetadata) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.CreativePage(pages_metadata).Render(context.Background(), c.Response().Writer)
	}
}

func readingListRenderer(pages_metadata *ds.PagesMetadata, books *map[string][]db.Book) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.ReadingListPage(pages_metadata, books).Render(context.Background(), c.Response().Writer)
	}
}

func blogPageRenderer(fname, url string, pages_metadata *ds.PagesMetadata, posts_metadata *ds.PostsMetadata) echo.HandlerFunc {
	return func(c echo.Context) error {
		return pages.BlogEntryPage(fname, url, pages_metadata, posts_metadata).Render(context.Background(), c.Response().Writer)
	}
}
