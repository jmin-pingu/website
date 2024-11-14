package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"sync"

	"mywebsite/pub/partials"

	"github.com/labstack/echo/v4"
)

var (
	count = 1
)

type Todo struct {
	list map[string]bool
	mu   sync.Mutex
}

var (
	todo = Todo{list: make(map[string]bool, 0), mu: sync.Mutex{}}
)

func getTodo(c echo.Context) error {
	vals := c.Request()
	log.Printf("updateTodo %v", vals)

	return c.NoContent(http.StatusOK)
}

func createTodo(c echo.Context) error {
	var b bytes.Buffer

	values, _ := c.FormParams()
	todo.mu.Lock()
	defer todo.mu.Unlock()
	log.Printf("createTodo: %v", values)
	entry := values["newTodo"][0]
	todo.list[entry] = false

	keys := make([]string, 0, len(todo.list))
	for v := range todo.list {
		keys = append(keys, v)
	}
	// Need to render the appropriate component
	component := partials.TodoEntryList(keys)
	component.Render(context.Background(), &b)
	log.Printf("createTodo: %v", b.String())
	return c.HTML(http.StatusOK, b.String())
}

func updateTodo(c echo.Context) error {
	vals, _ := c.FormParams()
	log.Printf("updateTodo %v", c.ParamValues())
	log.Printf("updateTodo %v", vals)

	return c.NoContent(http.StatusOK)
}

func deleteTodo(c echo.Context) error {
	// Need to make a way of deleting the item from todo
	vals, _ := c.FormParams()
	log.Printf("deleteTodo %v", vals)
	log.Printf("deleteTodo", c.ParamValues())
	log.Printf("deleteTodo", c.QueryParams())
	return c.NoContent(http.StatusOK)
}

// Element functionality
func logClick(c echo.Context) error {
	log.Println(count)
	count++
	return c.String(http.StatusOK, "Click Me")
}
