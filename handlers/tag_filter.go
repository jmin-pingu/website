package handlers

import (
	"bytes"
	"context"
	"log"
	"mywebsite/components/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
