package config

import (
	"fmt"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	DynamoDB        DynamoDBConfig `yaml:"dynamodb"`
	ShutdownTimeout string         `yaml:"shutdown_timeout"`
}

// DynamoDBConfig represents DynamoDB specific configuration
type DynamoDBConfig struct {
	Endpoint  string `yaml:"endpoint"`
	Region    string `yaml:"region"`
	TableName string `yaml:"table_name"`
	Timeout   string `yaml:"timeout"`
}

func LoadConfig() (*Config, error) {
	// Default configuration
	config := &Config{
		DynamoDB: DynamoDBConfig{
			Region:    "ap-northeast-1",
			TableName: "goto-dev-todo",
			Timeout:   "1s",
		},
		ShutdownTimeout: "5s",
	}

	// Override with environment variables if they exist
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		config.DynamoDB.Endpoint = endpoint
	}
	if region := os.Getenv("AWS_REGION"); region != "" {
		config.DynamoDB.Region = region
	}
	if table := os.Getenv("DYNAMODB_TABLE"); table != "" {
		config.DynamoDB.TableName = table
	}
	if timeout := os.Getenv("DYNAMODB_CONNECTION_TIMEOUT"); timeout != "" {
		config.DynamoDB.Timeout = timeout
	}
	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		config.ShutdownTimeout = timeout
	}

	// Validate durations
	if _, err := time.ParseDuration(config.DynamoDB.Timeout); err != nil {
		panic(fmt.Sprintf("Invalid format for DYNAMODB_CONNECTION_TIMEOUT: %v", err))
	}
	if _, err := time.ParseDuration(config.ShutdownTimeout); err != nil {
		panic(fmt.Sprintf("Invalid format for SHUTDOWN_TIMEOUT: %v", err))
	}

	return config, nil
}
