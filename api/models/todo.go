package models

// Todo represents a single TODO item in our application.
// The `json` tags are used by Gin to automatically bind JSON requests
// to this struct and to serialize this struct into JSON responses.
type Todo struct {
	ID        string `json:"id"`         // Unique identifier for the TODO item
	Title     string `json:"title" binding:"required"` // Title of the TODO item, required for creation/update
	Completed bool   `json:"completed"`  // Status of the TODO item (true if completed, false otherwise)
}
