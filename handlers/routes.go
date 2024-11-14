package handlers

import (
	"bytes"
	"context"
	"log"
	"mywebsite/pub/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Render pages
	SetupRenders(e)

	// set up static assets
	e.Static("/assets", "pub/assets")

	// BLOG
	// tag filtering functionality
	e.POST("/blog/", postFilter)
	// search functionality
	e.POST("/blog/search", getSearchInputs)
}

// getSearchInputs:
// TODO: add description
func getSearchInputs(c echo.Context) error {
	vals, _ := c.FormParams()
	log.Printf("getSearchInputs %v", vals["search"])

	var buf bytes.Buffer
	pages.BlogPosts(&POSTS_METADATA, DISPLAY_TAGS, vals["search"][0]).Render(context.Background(), &buf)
	return c.HTML(http.StatusOK, buf.String())
}

// postFilter:
// TODO: add description
func postFilter(c echo.Context) error {
	vals, _ := c.FormParams() // Get form parameters 1 - filter tag, 0 - do not filter tag
	for k, v := range vals {
		log.Printf("key %v, value %v", k, v)
		if v[0] == "1" {
			DISPLAY_TAGS.Add(k)
		} else {
			DISPLAY_TAGS.Remove(k)
		}
	}

	var buf bytes.Buffer
	pages.BlogPosts(&POSTS_METADATA, DISPLAY_TAGS, "").Render(context.Background(), &buf)
	return c.HTML(http.StatusOK, buf.String())
}
