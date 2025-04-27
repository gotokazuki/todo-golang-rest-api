package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/domain/entity"
	"github.com/gotokazuki/todo-golang-rest-api/app/todo/usecase/todo"
	"go.uber.org/zap"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	useCase *todo.TodoUseCase
	logger  *zap.Logger
}

// NewTodoHandler creates a new TodoHandler instance
func NewTodoHandler(useCase *todo.TodoUseCase, logger *zap.Logger) *TodoHandler {
	return &TodoHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// CreateTodo handles the creation of a new todo item
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var input entity.TodoCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid request body",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	createdTodo, err := h.useCase.CreateTodo(input)
	if err != nil {
		h.logger.Error("Failed to create todo",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create todo",
		})
		return
	}

	location := fmt.Sprintf("/todos/%s", createdTodo.ID.String())
	c.Header("Location", location)
	c.Status(http.StatusCreated)
}

// GetTodos handles retrieving all todo items
func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := h.useCase.GetTodos()
	if err != nil {
		h.logger.Error("Failed to get todos",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get todos",
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodo handles retrieving a specific todo item
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid todo ID",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid todo ID",
		})
		return
	}

	todo, err := h.useCase.GetTodo(id)
	if err != nil {
		h.logger.Error("Failed to get todo",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("id", id.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get todo",
		})
		return
	}

	if todo == nil {
		h.logger.Warn("Todo not found",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("id", id.String()),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// UpdateTodo handles updating an existing todo item
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid todo ID",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid todo ID",
		})
		return
	}

	var input entity.TodoUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid request body",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	updatedTodo, err := h.useCase.UpdateTodo(id, input)
	if err != nil {
		h.logger.Error("Failed to update todo",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("id", id.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update todo",
		})
		return
	}

	if updatedTodo == nil {
		h.logger.Warn("Todo not found",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("id", id.String()),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteTodo handles deleting a todo item
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid todo ID",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid todo ID",
		})
		return
	}

	err = h.useCase.DeleteTodo(id)
	if err != nil {
		h.logger.Error("Failed to delete todo",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("id", id.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete todo",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
