package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/domain/entity"
)

// TodoRepository defines the interface for todo data access
type TodoRepository interface {
	Create(todo *entity.Todo) (*entity.Todo, error)
	FindAll() ([]*entity.Todo, error)
	FindByID(id uuid.UUID) (*entity.Todo, error)
	Update(todo *entity.Todo) (*entity.Todo, error)
	Delete(id uuid.UUID) error
	GetClient() *dynamodb.Client
	GetTableName() string
}
