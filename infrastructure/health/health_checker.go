package health

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/config"
	"go.uber.org/zap"
)

// HealthStatus represents the status of a service
type HealthStatus string

const (
	StatusOK   HealthStatus = "ok"
	StatusFail HealthStatus = "fail"
)

// ServiceHealth represents the health status of a service
type ServiceHealth struct {
	Status  HealthStatus `json:"status"`
	Message string       `json:"message,omitempty"`
}

// HealthResponse represents the overall health check response
type HealthResponse struct {
	Status     string                   `json:"status"`
	Components map[string]ServiceHealth `json:"components"`
}

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	Check(ctx context.Context) HealthResponse
}

// DynamoDBHealthChecker implements HealthChecker for DynamoDB
type DynamoDBHealthChecker struct {
	client    *dynamodb.Client
	tableName string
	timeout   time.Duration
	logger    *zap.Logger
}

// NewDynamoDBHealthChecker creates a new DynamoDBHealthChecker
func NewDynamoDBHealthChecker(client *dynamodb.Client, tableName string, cfg *config.Config, logger *zap.Logger) *DynamoDBHealthChecker {

	// Parse DynamoDB timeout from string to time.Duration
	timeout, err := time.ParseDuration(cfg.DynamoDB.Timeout)
	if err != nil {
		panic(fmt.Sprintf("Invalid DynamoDB timeout format: %v", err))
	}

	return &DynamoDBHealthChecker{
		client:    client,
		tableName: tableName,
		timeout:   timeout,
		logger:    logger,
	}
}

// Check performs the health check for DynamoDB
func (h *DynamoDBHealthChecker) Check(ctx context.Context) HealthResponse {
	response := HealthResponse{
		Status:     string(StatusOK),
		Components: make(map[string]ServiceHealth),
	}

	// Check DynamoDB
	if err := h.checkDynamoDB(ctx); err != nil {
		h.logger.Error("DynamoDB health check failed",
			zap.Error(err),
			zap.String("component", "dynamodb"),
		)
		response.Status = string(StatusFail)
		response.Components["dynamodb"] = ServiceHealth{
			Status:  StatusFail,
			Message: "Failed to connect to todos table",
		}
		return response
	}

	response.Components["dynamodb"] = ServiceHealth{
		Status: StatusOK,
	}
	return response
}

// checkDynamoDB performs the actual DynamoDB health check
func (h *DynamoDBHealthChecker) checkDynamoDB(ctx context.Context) error {
	// Set a timeout for the health check
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	_, err := h.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &h.tableName,
	})
	return err
}
