package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/domain/entity"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/domain/repository"
)

// TodoUseCase handles the business logic for todo operations
type TodoUseCase struct {
	repo repository.TodoRepository
}

// NewTodoUseCase creates a new TodoUseCase instance
func NewTodoUseCase(repo repository.TodoRepository) *TodoUseCase {
	return &TodoUseCase{repo: repo}
}

// CreateTodo creates a new todo item
func (u *TodoUseCase) CreateTodo(input entity.TodoCreate) (*entity.Todo, error) {
	now := time.Now()
	todo := &entity.Todo{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return u.repo.Create(todo)
}

// GetTodos retrieves all todo items
func (u *TodoUseCase) GetTodos() ([]*entity.Todo, error) {
	return u.repo.FindAll()
}

// GetTodo retrieves a todo item by ID
func (u *TodoUseCase) GetTodo(id uuid.UUID) (*entity.Todo, error) {
	return u.repo.FindByID(id)
}

// UpdateTodo updates an existing todo item
func (u *TodoUseCase) UpdateTodo(id uuid.UUID, input entity.TodoUpdate) (*entity.Todo, error) {
	todo, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		todo.Title = input.Title
	}
	if input.Description != "" {
		todo.Description = input.Description
	}
	todo.Completed = input.Completed
	todo.UpdatedAt = time.Now()

	return u.repo.Update(todo)
}

// DeleteTodo deletes a todo item
func (u *TodoUseCase) DeleteTodo(id uuid.UUID) error {
	return u.repo.Delete(id)
}
