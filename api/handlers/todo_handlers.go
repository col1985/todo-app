package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // For generating unique IDs

	"my-todo-app-api/models"  // Import the Todo model
	"my-todo-app-api/storage" // Import the in-memory storage
)

// CreateTodo handles the creation of a new TODO item.
// It expects a JSON payload with at least a "title".
// The "completed" field is optional and defaults to false.
func CreateTodo(c *gin.Context) {
	var newTodo models.Todo

	// Bind the incoming JSON request body to the newTodo struct.
	// Gin automatically handles JSON parsing and validation based on struct tags.
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		// If binding fails (e.g., "title" is missing), return a 400 Bad Request.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new unique ID for the TODO item.
	newTodo.ID = uuid.New().String()
	// Ensure completed is false by default if not provided or invalid.
	// If the client sends `completed: null` or `completed: 0`, Gin's binding might set it to false.
	// If the client sends `completed: true`, it will be honored.

	// Store the new TODO item in our in-memory map.
	storage.Todos[newTodo.ID] = newTodo

	// Return the newly created TODO item with a 201 Created status.
	c.JSON(http.StatusCreated, newTodo)
}

// GetTodos retrieves all TODO items from storage.
func GetTodos(c *gin.Context) {
	// Convert the map of TODOs into a slice of TODOs for JSON serialization.
	// This is necessary because Gin's JSON marshaller expects a slice or array for a list.
	todos := []models.Todo{}
	for _, todo := range storage.Todos {
		todos = append(todos, todo)
	}

	// Return the slice of TODO items with a 200 OK status.
	c.JSON(http.StatusOK, todos)
}

// GetTodoByID retrieves a single TODO item by its ID.
// The ID is expected as a path parameter, e.g., /todos/:id.
func GetTodoByID(c *gin.Context) {
	// Get the ID from the URL path parameters.
	id := c.Param("id")

	// Look up the TODO item in our storage map.
	todo, exists := storage.Todos[id]
	if !exists {
		// If the TODO item is not found, return a 404 Not Found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// Return the found TODO item with a 200 OK status.
	c.JSON(http.StatusOK, todo)
}

// UpdateTodo updates an existing TODO item by its ID.
// It expects the ID as a path parameter and a JSON payload for updates.
func UpdateTodo(c *gin.Context) {
	// Get the ID from the URL path parameters.
	id := c.Param("id")

	// Check if the TODO item exists before attempting to update.
	existingTodo, exists := storage.Todos[id]
	if !exists {
		// If not found, return a 404 Not Found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var updatedTodo models.Todo
	// Bind the incoming JSON request body to the updatedTodo struct.
	// This will overwrite fields present in the JSON.
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		// If binding fails, return a 400 Bad Request.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the fields of the existing TODO item.
	// We only update Title and Completed if they are provided in the request.
	// The ID should not be changed during an update.
	if updatedTodo.Title != "" {
		existingTodo.Title = updatedTodo.Title
	}
	// The `Completed` field is a boolean, so we need to check if it was explicitly provided.
	// Gin's ShouldBindJSON will set a boolean to its zero value (false) if not present.
	// A more robust solution for partial updates would involve custom unmarshalling or checking
	// if the field was truly present in the JSON. For simplicity, we'll update if the title is non-empty
	// and always update `Completed` if it's part of the bound struct.
	// If you want to allow setting `Completed` to false explicitly, the current binding is fine.
	// If you want to only update `Completed` if it's explicitly sent, you'd need a pointer `*bool` in the model
	// or a custom unmarshal logic. For this basic example, we'll assume if `completed` is in the payload, it updates.
	existingTodo.Completed = updatedTodo.Completed


	// Store the updated TODO item back into the map.
	storage.Todos[id] = existingTodo

	// Return the updated TODO item with a 200 OK status.
	c.JSON(http.StatusOK, existingTodo)
}

// DeleteTodo deletes a TODO item by its ID.
// The ID is expected as a path parameter.
func DeleteTodo(c *gin.Context) {
	// Get the ID from the URL path parameters.
	id := c.Param("id")

	// Check if the TODO item exists before attempting to delete.
	_, exists := storage.Todos[id]
	if !exists {
		// If not found, return a 404 Not Found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// Delete the TODO item from the map.
	delete(storage.Todos, id)

	// Return a 204 No Content status to indicate successful deletion with no response body.
	c.Status(http.StatusNoContent)
}
