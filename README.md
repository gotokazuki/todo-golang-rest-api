# todo-golang-rest-api

A RESTful API for managing TODO items, built with Go.

## Features

- Create, read, update, and delete TODO items.
- DynamoDB integration for data persistence.
- Health check endpoint for monitoring service status.

## Requirements

- Go 1.20 or later
- Docker (for running dependencies like DynamoDB Local or LocalStack)

## Configuration

The application is configured using environment variables. Below are the available variables and their default values:

| Environment Variable          | Description                     | Default Value           |
| ----------------------------- | ------------------------------- | ----------------------- |
| `DYNAMODB_ENDPOINT`           | DynamoDB endpoint URL           | `http://localhost:4566` |
| `AWS_REGION`                  | AWS region                      | `ap-northeast-1`        |
| `DYNAMODB_TABLE`              | DynamoDB table name             | `goto-dev-todo`         |
| `DYNAMODB_CONNECTION_TIMEOUT` | Timeout for DynamoDB operations | `1s`                    |
| `SHUTDOWN_TIMEOUT`            | Timeout for graceful shutdown   | `5s`                    |

## Running the Application

### Using Docker

1. Build the Docker image:

```shell
make build
```

2. Run the container:

```shell
docker compose up -d
```

### Running Locally with Amazon DynamoDB

To run the application with Amazon DynamoDB, ensure that your AWS credentials are properly configured.

1. Install dependencies:

```shell
go mod tidy
```

2. Run the application:

```shell
go run main.go
```

### Running Locally with LocalStack on Docker

LocalStack is a fully functional local AWS cloud stack. Use it to mock AWS services like DynamoDB.

1. Start localstack:

```shell
docker compose up -d localstack
```

2. Install dependencies:

```shell
go mod tidy
```

3. Run the application with variables:

```shell
DYNAMODB_ENDPOINT=http://localhost:4566 go run main.go
```

## Endpoints

| Method | Endpoint       | Description              |
|--------|----------------|--------------------------|
| GET    | `/todos`       | Get all TODO items       |
| POST   | `/todos`       | Create a new TODO item   |
| GET    | `/todos/{id}`  | Get a TODO item by ID    |
| PUT    | `/todos/{id}`  | Update a TODO item by ID |
| DELETE | `/todos/{id}`  | Delete a TODO item by ID |
| GET    | `/health`      | Health check endpoint    |

For more details, see [docs/openapi.yaml](./docs/openapi.yaml).

## Health Check

The `/health` endpoint provides the status of the service and its dependencies (e.g., DynamoDB).

### Example Request

```shell
curl -s localhost:8080/health | jq .
```

### Example Response

```json
{
  "status": "ok",
  "components": {
    "dynamodb": {
      "status": "ok"
    }
  }
}
```

## Test

A simple shell script `test.sh` is provided to verify the functionality of the application.
This script performs basic tests, such as checking the health endpoint and creating a TODO item.

### Usage

1. Ensure the application is running locally or in Docker.
2. Run the test script:

```shell
./test.sh
```

### Example Output

```shell
Checking health endpoint...
{"status":"ok","components":{"dynamodb":{"status":"ok"}}}
Health check passed!
Creating a TODO item...
TODO item created successfully! The created item location: /todos/a22b5f8a-c698-4f48-ba76-11e9e9efebdb
Fetch all TODO items...
[{"ID":"a22b5f8a-c698-4f48-ba76-11e9e9efebdb","Title":"Sample Todo","Description":"This is a test todo","Completed":false,"CreatedAt":"2025-04-27T06:39:07Z","UpdatedAt":"2025-04-27T06:39:07Z"}]
Updating a TODO item...
TODO item updated successfully!
Fetching an updated TODO item...
{"ID":"a22b5f8a-c698-4f48-ba76-11e9e9efebdb","Title":"Updated Todo","Description":"This is an updated test todo","Completed":true,"CreatedAt":"2025-04-27T06:39:07Z","UpdatedAt":"2025-04-27T06:39:07Z"}
Deleting a TODO item...
TODO item deleted successfully!
All tests completed successfully!
```
