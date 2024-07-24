package handlers

import (
	"bytes"
	"context"
	"log"
	"mywebsite/components/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getSearchInputs(c echo.Context) error {
	vals, _ := c.FormParams()
	log.Printf("getSearchInputs %v", vals["search"])

	var buf bytes.Buffer
	pages.BlogPosts(&POSTS_METADATA, DISPLAY_TAGS, vals["search"][0]).Render(context.Background(), &buf)
	return c.HTML(http.StatusOK, buf.String())
}
