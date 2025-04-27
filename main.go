package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/config"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/dynamodb"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/health"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/logger"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/infrastructure/middleware"
	todohttp "github.com/gotokazuki/todo-golang-rest-api/app/todo/interface/http"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/usecase/todo"
	"go.uber.org/zap"
)

// main is the entry point of the application.
// It initializes the repository, handler, and sets up the Gin router with middleware and routes.
func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	cfg, err := config.LoadConfig()

	// Initialize repository
	repo := dynamodb.NewTodoRepository(cfg)

	// Initialize use case
	useCase := todo.NewTodoUseCase(repo)

	// Initialize handler
	handler := todohttp.NewTodoHandler(useCase, log)

	// Initialize health checker
	healthChecker := health.NewDynamoDBHealthChecker(repo.GetClient(), repo.GetTableName(), cfg, log)

	// Create Gin router without default middleware
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(log))
	r.Use(middleware.ErrorHandlerMiddleware(log))

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	r.GET("/health", func(c *gin.Context) {
		response := healthChecker.Check(c.Request.Context())
		c.JSON(getStatusCode(response.Status), response)
	})
	r.GET("/todos", handler.GetTodos)
	r.POST("/todos", handler.CreateTodo)
	r.GET("/todos/:id", handler.GetTodo)
	r.PATCH("/todos/:id", handler.UpdateTodo)
	r.DELETE("/todos/:id", handler.DeleteTodo)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	timeout, err := time.ParseDuration(cfg.ShutdownTimeout)
	if err != nil {
		panic(fmt.Sprintf("Invalid shutdown timeout format: %v", err))
	}

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}

// getStatusCode returns the appropriate HTTP status code based on the health status
func getStatusCode(status string) int {
	switch status {
	case string(health.StatusOK):
		return 200
	case string(health.StatusFail):
		return 500
	default:
		return 500
	}
}
