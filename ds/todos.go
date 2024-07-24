package ds

import (
	"strings"

	"github.com/google/uuid"
)

type Todos []*Todo

func NewTodos() *Todos {
	return &Todos{}
}

func (todos *Todos) Add(description string) *Todo {
	todo := NewTodo(description)
	*todos = append(*todos, todo)
	return todo
}

func (todos *Todos) Delete(id uuid.UUID) *Todo {
	i := todos.index(id)
	*todos = append((*todos)[:i], (*todos)[i+1:]...)
	return (*todos)[i]
}

func (todos *Todos) Get(id uuid.UUID) *Todo {
	i := todos.index(id)
	if i == -1 {
		return nil
	}
	return (*todos)[i]
}

func (todos *Todos) Search(search string) *Todos {
	l := NewTodos()
	for _, todo := range *todos {
		if strings.Contains(todo.Description, search) {
			*l = append(*l, todo)
		}
	}
	return l
}

func (todos *Todos) Update(id uuid.UUID, completed bool, description string) *Todo {
	i := todos.index(id)
	if i == -1 {
		return nil
	}
	(*todos)[i].Update(completed, description)
	return (*todos)[i]
}

func (todos *Todos) index(id uuid.UUID) int {
	for i, todo := range *todos {
		if todo.ID == id {
			return i
		}
	}
	return -1
}
