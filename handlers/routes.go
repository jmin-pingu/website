package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Render pages
	SetupRenders(e)

	// BLOG
	// tag filtering functionality
	e.POST("/blog/", postFilter)
	// search functionality
	e.POST("/blog/search", getSearchInputs)
}
