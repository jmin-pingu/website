package main

import (
	"mywebsite/handlers"

	//"github.com/kaleocheng/goldmark"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Setup routes
	handlers.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
