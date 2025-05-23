Go Gin TODO API
This is a simple RESTful TODO API built with Go and the Gin web framework. It includes in-memory storage, CORS support, and Prometheus metrics exposure.

Table of Contents
Features

Prerequisites

Getting Started

Build and Run Locally

Build and Run with Docker

API Endpoints

Health Check

Prometheus Metrics

Create TODO

Get All TODOs

Get TODO by ID

Update TODO

Delete TODO

Testing

Project Structure

Features
CRUD Operations: Create, Read, Update, and Delete TODO items.

Gin Framework: Fast and lightweight web framework for Go.

In-Memory Storage: Simple data storage (can be easily replaced with a database).

CORS Support: Configured to allow requests from http://localhost:3000 (for a React frontend).

Prometheus Metrics: Exposes application metrics at /metrics endpoint.

Unit Tests: Basic unit tests for handler functions.

Dockerfile: For easy containerization.

Prerequisites
Before you begin, ensure you have the following installed:

Go: Version 1.18 or higher.

Git: For cloning the repository.

Docker (Optional, for containerized deployment).

Getting Started
Follow these steps to get the API up and running.

Build and Run Locally
Clone the repository:

git clone <repository-url> # Replace with your actual repository URL
cd todo-app

Initialize Go Modules and Download Dependencies:

go mod tidy

Run the application:

go run main.go

The API will start on http://localhost:8080.

Build and Run with Docker
Build the Docker image:

docker build -t go-gin-todo-api .

Run the Docker container:

docker run -p 8080:8080 go-gin-todo-api

The API will be accessible at http://localhost:8080.

API Endpoints
The API base URL is http://localhost:8080/api/v1.

Health Check
Checks the health of the API.

Endpoint: GET /health

Response:

{
    "status": "ok"
}

Curl Command:

curl http://localhost:8080/health

Prometheus Metrics
Exposes application metrics in Prometheus format.

Endpoint: GET /metrics

Response: Text-based Prometheus metrics.

Curl Command:

curl http://localhost:8080/metrics

Create TODO
Creates a new TODO item.

Endpoint: POST /api/v1/todos

Request Body (JSON):

{
    "title": "My new todo item",
    "completed": false
}

title is required.

completed is optional (defaults to false).

Response (201 Created):

{
    "id": "generated-uuid",
    "title": "My new todo item",
    "completed": false
}

Curl Command:

curl -X POST -H "Content-Type: application/json" -d '{"title": "Learn Go Gin", "completed": false}' http://localhost:8080/api/v1/todos

Get All TODOs
Retrieves all TODO items.

Endpoint: GET /api/v1/todos

Response (200 OK):

[
    {
        "id": "uuid-1",
        "title": "Learn Go Gin",
        "completed": false
    },
    {
        "id": "uuid-2",
        "title": "Build a REST API",
        "completed": true
    }
]

Curl Command:

curl http://localhost:8080/api/v1/todos

Get TODO by ID
Retrieves a single TODO item by its ID.

Endpoint: GET /api/v1/todos/:id

Response (200 OK):

{
    "id": "uuid-of-todo",
    "title": "Specific Todo Item",
    "completed": false
}

Response (404 Not Found):

{
    "error": "Todo not found"
}

Curl Command (replace [ID] with an actual TODO ID):

curl http://localhost:8080/api/v1/todos/[ID]

Update TODO
Updates an existing TODO item by its ID.

Endpoint: PUT /api/v1/todos/:id

Request Body (JSON):

{
    "title": "Updated todo title",
    "completed": true
}

Provide only the fields you want to update. title is required if present.

Response (200 OK):

{
    "id": "uuid-of-todo",
    "title": "Updated todo title",
    "completed": true
}

Response (404 Not Found):

{
    "error": "Todo not found"
}

Curl Command (replace [ID] with an actual TODO ID):

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"title": "Updated title", "completed": true}' http://localhost:8080/api/v1/todos/[ID]
```

Delete TODO
Deletes a TODO item by its ID.

Endpoint: DELETE /api/v1/todos/:id

Response (204 No Content): Successful deletion.

Response (404 Not Found):

{
    "error": "Todo not found"
}

Curl Command (replace [ID] with an actual TODO ID):

curl -X DELETE http://localhost:8080/api/v1/todos/[ID]

Testing
To run the unit tests for the API handlers:

go test ./...

Project Structure
todo-app/
├── main.go               # Main application entry point, router setup, CORS, Prometheus
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Dockerfile            # Docker configuration for containerization
├── handlers/
│   └── todo_handlers.go  # API handler functions (Create, Get, Update, Delete)
│   └── todo_handlers_test.go # Unit tests for handlers
├── models/
│   └── todo.go           # Defines the Todo struct
└── storage/
    └── inmemory.go       # In-memory storage for Todo items
