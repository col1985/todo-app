package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert" // Popular assertion library

	"my-todo-app-api/handlers"
	"my-todo-app-api/models"
	"my-todo-app-api/storage"
)

// setupRouter sets up a Gin router for testing purposes.
// It initializes the in-memory storage and registers the API routes.
func setupRouter() *gin.Engine {
	r := gin.Default()
	storage.Todos = make(map[string]models.Todo) // Clear storage for each test run

	v1 := r.Group("/api/v1")
	{
		v1.POST("/todos", handlers.CreateTodo)
		v1.GET("/todos", handlers.GetTodos)
		v1.GET("/todos/:id", handlers.GetTodoByID)
		v1.PUT("/todos/:id", handlers.UpdateTodo)
		v1.DELETE("/todos/:id", handlers.DeleteTodo)
	}
	return r
}

// TestCreateTodo tests the CreateTodo handler.
func TestCreateTodo(t *testing.T) {
	router := setupRouter()

	// Define a new todo item to be created
	newTodo := models.Todo{
		Title: "Test Create Todo",
	}
	jsonValue, _ := json.Marshal(newTodo)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the HTTP response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(w, req)

	// Assert the HTTP status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response body
	var createdTodo models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &createdTodo)
	assert.NoError(t, err)

	// Assert the fields of the created todo
	assert.NotEmpty(t, createdTodo.ID)
	assert.Equal(t, newTodo.Title, createdTodo.Title)
	assert.False(t, createdTodo.Completed) // Should be false by default

	// Verify it's in storage
	_, exists := storage.Todos[createdTodo.ID]
	assert.True(t, exists)
}

// TestCreateTodo_InvalidInput tests CreateTodo with invalid input (missing title).
func TestCreateTodo_InvalidInput(t *testing.T) {
	router := setupRouter()

	// Invalid todo item (missing title, which is required)
	invalidTodo := models.Todo{
		Completed: true,
	}
	jsonValue, _ := json.Marshal(invalidTodo)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Key: 'Todo.Title' Error:Field validation for 'Title' failed on the 'required' tag")
}

// TestGetTodos tests the GetTodos handler.
func TestGetTodos(t *testing.T) {
	router := setupRouter()

	// Pre-populate storage with some todos
	id1 := uuid.New().String()
	storage.Todos[id1] = models.Todo{ID: id1, Title: "Todo 1", Completed: false}
	id2 := uuid.New().String()
	storage.Todos[id2] = models.Todo{ID: id2, Title: "Todo 2", Completed: true}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var todos []models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &todos)
	assert.NoError(t, err)
	assert.Len(t, todos, 2) // Should contain the two pre-populated todos
}

// TestGetTodoByID tests the GetTodoByID handler for an existing todo.
func TestGetTodoByID(t *testing.T) {
	router := setupRouter()

	// Pre-populate a specific todo
	id := uuid.New().String()
	expectedTodo := models.Todo{ID: id, Title: "Specific Todo", Completed: false}
	storage.Todos[id] = expectedTodo

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var todo models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &todo)
	assert.NoError(t, err)
	assert.Equal(t, expectedTodo.ID, todo.ID)
	assert.Equal(t, expectedTodo.Title, todo.Title)
	assert.Equal(t, expectedTodo.Completed, todo.Completed)
}

// TestGetTodoByID_NotFound tests GetTodoByID for a non-existent todo.
func TestGetTodoByID_NotFound(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/todos/nonexistent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Todo not found")
}

// TestUpdateTodo tests the UpdateTodo handler.
func TestUpdateTodo(t *testing.T) {
	router := setupRouter()

	// Pre-populate a todo to be updated
	id := uuid.New().String()
	originalTodo := models.Todo{ID: id, Title: "Original Title", Completed: false}
	storage.Todos[id] = originalTodo

	// Define update payload
	updatedTodo := models.Todo{
		Title:     "Updated Title",
		Completed: true,
	}
	jsonValue, _ := json.Marshal(updatedTodo)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/todos/"+id, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseTodo models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &responseTodo)
	assert.NoError(t, err)

	// Assert the updated fields
	assert.Equal(t, id, responseTodo.ID) // ID should remain the same
	assert.Equal(t, updatedTodo.Title, responseTodo.Title)
	assert.Equal(t, updatedTodo.Completed, responseTodo.Completed)

	// Verify it's updated in storage
	storedTodo, _ := storage.Todos[id]
	assert.Equal(t, updatedTodo.Title, storedTodo.Title)
	assert.Equal(t, updatedTodo.Completed, storedTodo.Completed)
}

// TestUpdateTodo_NotFound tests UpdateTodo for a non-existent todo.
func TestUpdateTodo_NotFound(t *testing.T) {
	router := setupRouter()

	updatePayload := models.Todo{Title: "Non Existent", Completed: true}
	jsonValue, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/todos/nonexistent-id", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Todo not found")
}

// TestUpdateTodo_InvalidInput tests UpdateTodo with invalid JSON input.
func TestUpdateTodo_InvalidInput(t *testing.T) {
	router := setupRouter()

	id := uuid.New().String()
	storage.Todos[id] = models.Todo{ID: id, Title: "Original", Completed: false}

	// Send invalid JSON (e.g., not a valid Todo struct)
	invalidJson := `{"title": 123}` // Title should be string
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/todos/"+id, bytes.NewBufferString(invalidJson))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "json: cannot unmarshal number into Go struct field Todo.title of type string")
}


// TestDeleteTodo tests the DeleteTodo handler.
func TestDeleteTodo(t *testing.T) {
	router := setupRouter()

	// Pre-populate a todo to be deleted
	id := uuid.New().String()
	storage.Todos[id] = models.Todo{ID: id, Title: "Todo to Delete", Completed: false}

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/todos/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code) // 204 No Content for successful deletion

	// Verify it's deleted from storage
	_, exists := storage.Todos[id]
	assert.False(t, exists)
}

// TestDeleteTodo_NotFound tests DeleteTodo for a non-existent todo.
func TestDeleteTodo_NotFound(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/todos/nonexistent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Todo not found")
}
