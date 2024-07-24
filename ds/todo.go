package ds

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID
	Description string
	Completed   bool
	CreatedAt   time.Time
}

func NewTodo(description string) *Todo {
	return &Todo{
		ID:          uuid.New(),
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}

func (todo *Todo) Update(completed bool, description string) {
	todo.Completed = completed
	todo.Description = description
}
