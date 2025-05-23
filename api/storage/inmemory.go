package storage

import "my-todo-app-api/models"

// Todos is a global map that serves as our in-memory storage for Todo items.
// The key is the Todo ID (string), and the value is the Todo struct.
// In a production application, this would be replaced by a database connection.
var Todos map[string]models.Todo
