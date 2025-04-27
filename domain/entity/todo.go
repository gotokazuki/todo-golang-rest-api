package entity

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item in the domain
type Todo struct {
	ID          uuid.UUID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TodoCreate represents the data needed to create a new todo
type TodoCreate struct {
	Title       string
	Description string
}

// TodoUpdate represents the data needed to update an existing todo
type TodoUpdate struct {
	Title       string
	Description string
	Completed   bool
}
