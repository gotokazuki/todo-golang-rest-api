package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/domain/entity"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/domain/repository"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/config"
)

// TodoRepository implements the repository.TodoRepository interface for DynamoDB
type TodoRepository struct {
	client  *dynamodb.Client
	table   string
	timeout time.Duration
}

// NewTodoRepository creates a new TodoRepository instance
func NewTodoRepository(cfg *config.Config) repository.TodoRepository {
	var awsCfg aws.Config
	var err error

	if cfg.DynamoDB.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           cfg.DynamoDB.Endpoint,
				SigningRegion: region,
			}, nil
		})

		awsCfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithEndpointResolverWithOptions(customResolver),
			awsconfig.WithRegion(cfg.DynamoDB.Region),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
		)
	} else {
		awsCfg, err = awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithRegion(cfg.DynamoDB.Region),
		)
	}

	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(awsCfg)

	// Parse DynamoDB timeout from string to time.Duration
	timeout, err := time.ParseDuration(cfg.DynamoDB.Timeout)
	if err != nil {
		panic(fmt.Sprintf("Invalid DynamoDB timeout format: %v", err))
	}

	return &TodoRepository{
		client:  client,
		table:   cfg.DynamoDB.TableName,
		timeout: timeout,
	}
}

// withTimeout creates a context with the configured timeout
func (r *TodoRepository) withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.timeout)
}

// Create saves a new todo item to DynamoDB
func (r *TodoRepository) Create(todo *entity.Todo) (*entity.Todo, error) {
	ctx, cancel := r.withTimeout()
	defer cancel()

	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: todo.ID.String()},
		"title":       &types.AttributeValueMemberS{Value: todo.Title},
		"description": &types.AttributeValueMemberS{Value: todo.Description},
		"completed":   &types.AttributeValueMemberBOOL{Value: todo.Completed},
		"created_at":  &types.AttributeValueMemberS{Value: todo.CreatedAt.Format("2006-01-02T15:04:05Z")},
		"updated_at":  &types.AttributeValueMemberS{Value: todo.UpdatedAt.Format("2006-01-02T15:04:05Z")},
	}

	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      item,
	})
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// FindAll retrieves all todo items from DynamoDB
func (r *TodoRepository) FindAll() ([]*entity.Todo, error) {
	ctx, cancel := r.withTimeout()
	defer cancel()

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.table),
	})
	if err != nil {
		return nil, err
	}

	todos := make([]*entity.Todo, 0, len(result.Items))
	for _, item := range result.Items {
		todo, err := r.unmarshalTodo(item)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// FindByID retrieves a todo item by its ID from DynamoDB
func (r *TodoRepository) FindByID(id uuid.UUID) (*entity.Todo, error) {
	ctx, cancel := r.withTimeout()
	defer cancel()

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	return r.unmarshalTodo(result.Item)
}

// Update saves changes to an existing todo item in DynamoDB
func (r *TodoRepository) Update(todo *entity.Todo) (*entity.Todo, error) {
	ctx, cancel := r.withTimeout()
	defer cancel()

	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: todo.ID.String()},
			"title":       &types.AttributeValueMemberS{Value: todo.Title},
			"description": &types.AttributeValueMemberS{Value: todo.Description},
			"completed":   &types.AttributeValueMemberBOOL{Value: todo.Completed},
			"created_at":  &types.AttributeValueMemberS{Value: todo.CreatedAt.Format("2006-01-02T15:04:05Z")},
			"updated_at":  &types.AttributeValueMemberS{Value: todo.UpdatedAt.Format("2006-01-02T15:04:05Z")},
		},
	})
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// Delete removes a todo item from DynamoDB by its ID
func (r *TodoRepository) Delete(id uuid.UUID) error {
	ctx, cancel := r.withTimeout()
	defer cancel()

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	return err
}

// unmarshalTodo converts a DynamoDB item to a Todo entity
func (r *TodoRepository) unmarshalTodo(item map[string]types.AttributeValue) (*entity.Todo, error) {
	idStr, ok := item["id"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, errors.New("invalid id type")
	}

	id, err := uuid.Parse(idStr.Value)
	if err != nil {
		return nil, err
	}

	title, ok := item["title"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, errors.New("invalid title type")
	}

	description, ok := item["description"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, errors.New("invalid description type")
	}

	completed, ok := item["completed"].(*types.AttributeValueMemberBOOL)
	if !ok {
		return nil, errors.New("invalid completed type")
	}

	createdAtStr, ok := item["created_at"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, errors.New("invalid created_at type")
	}

	createdAt, err := time.Parse("2006-01-02T15:04:05Z", createdAtStr.Value)
	if err != nil {
		return nil, err
	}

	updatedAtStr, ok := item["updated_at"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, errors.New("invalid updated_at type")
	}

	updatedAt, err := time.Parse("2006-01-02T15:04:05Z", updatedAtStr.Value)
	if err != nil {
		return nil, err
	}

	return &entity.Todo{
		ID:          id,
		Title:       title.Value,
		Description: description.Value,
		Completed:   completed.Value,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

// GetClient returns the DynamoDB client
func (r *TodoRepository) GetClient() *dynamodb.Client {
	return r.client
}

// GetTableName returns the table name
func (r *TodoRepository) GetTableName() string {
	return r.table
}
