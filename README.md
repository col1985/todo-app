# Go Gin API & React UI TODO Application

This repository contains a full-stack TODO application, consisting of a Go Gin RESTful API backend and a ReactJS single-page application (SPA) frontend.

Table of Contents
Project Overview

Features

Prerequisites

Getting Started

Project Structure

Backend (Go Gin API)

Build and Run Locally

Build and Run with Docker

Frontend (ReactJS UI)

Build and Run Locally

Build and Run with Docker

Running Both Applications

API Endpoints

Health Check

Prometheus Metrics

Create TODO

Get All TODOs

Get TODO by ID

Update TODO

Delete TODO

Testing

## Project Overview

This application demonstrates a simple TODO list functionality. The backend is a RESTful API written in Go using the Gin framework, providing CRUD operations for TODO items. It includes in-memory storage, CORS support for cross-origin requests, and a Prometheus metrics endpoint for monitoring. The frontend is a ReactJS single-page application that consumes this API, providing an interactive user interface to manage TODOs.

### Features

#### API (Go Gin App)

- CRUD Operations: Create, Read, Update, and Delete TODO items.
- Gin Framework: Fast and lightweight web framework for Go.
- In-Memory Storage: Simple data storage (can be easily replaced with a database).
- CORS Support: Configured to allow requests from `http://localhost:3000`.
- Prometheus Metrics: Exposes application metrics at /metrics endpoint.
- Unit Tests: Basic unit tests for handler functions.
- Containerfile: For easy containerization.
- Interactive Documentation: Serves an interactive HTML documentation page at `/docs/index.html`.

#### UI (ReactJS UI)

- Interactive TODO List: User-friendly interface for managing TODOs.
- Real-time Updates: Reflects changes in the backend API.
- Add, Toggle, Edit, Delete: Full functionality to interact with TODO items.
- Responsive Design: Built with Tailwind CSS for adaptability across devices.
- Containerization: Production-ready Containerfile using Nginx for serving.

### Prerequisites

Before you begin, ensure you have the following installed:

- Go: Version 1.23 or higher.
- Node.js & npm (or Yarn): Version 16 or higher for the React application.
- Git: For cloning the repository.
- Podman (Optional, but recommended for consistent environments).

### Getting Started

#### Project Structure

```bash
├── api/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── Containerfile
│   ├── handlers/
│   │   ├── todo_handlers.go
│   │   └── todo_handlers_test.go
│   ├── models/
│   │   └── todo.go
│   ├── storage/
│   │   └── inmemory.go
│   └── static/
│       ├── test.html
│       └── index.html
└── ui/
    ├── public/
    ├── src/
    ├── package.json
    ├── package-lock.json
    ├── Containerfile
    └── nginx.conf
```

### Backend (API)

Navigate to the `api` directory:

```bash
cd api
```

#### Build and Run Locally

Initialize Go Modules and Download Dependencies:

```bash
go mod tidy
```

Run the application:

```bash
go run main.go
```

The API will start on http://localhost:8080.

The interactive API documentation will be available at http://localhost:8080/ (redirects to /docs/index.html).

#### Build and Run with Podman

Build the Docker image:

```bash
podman build -t todo-api:latest .
```

Run the container:

```bash
podman run -p 8080:8080 todo-api:latest
```

The API will be accessible at http://localhost:8080.

### Frontend (UI)

Navigate to the `ui` directory:

```bash
cd ui
```

#### Build and Run Locally (UI)

Install Node.js Dependencies:

```bash
npm i # or yarn install
```

Start the dev server:

```bash
npm start # or yarn start
```

The UI will open in your browser at http://localhost:3000.

#### Build and Run with Podman

Ensure `nginx.conf` is in the `ui` root directory.

Build the image:

```bash
podman build -t todo-ui:latest .
```

Run the container:

```bash
podman run -p 3000:80 todo-ui:latest
```

The UI will be accessible at http://localhost:3000.

### Running Both Applications

For the full application experience, you need to run both the API backend and the UI frontend.

#### Start the Backend

Open a terminal, navigate to the todo-api directory, and run:

```bash
go run main.go
# or if using Docker:
podman run -p 8080:8080 todo-api:latest
```

### Start the Frontend

Open a separate terminal, navigate to the `ui` directory, and run:

```bash
npm start
# or if using Podman:
podman run -p 3000:80 todo-ui:latest
```

Once both are running, open your browser to http://localhost:3000 to interact with the TODO application.

### API Endpoints

The API base URL is http://localhost:8080/api/v1.

#### Health Check

Checks the health of the API.

**Endpoint:** *GET* `/health`

Response:

```json
{
    "status": "ok"
}
```

Curl Command:

```bash
curl http://localhost:8080/health
```

#### Prometheus Metrics

Exposes application metrics in Prometheus format.

**Endpoint**: *GET* `/metrics`

Response: Text-based Prometheus metrics.

Curl Command:

```bash
curl http://localhost:8080/metrics
```

#### Create a TODO

Creates a new TODO item.

**Endpoint**: *POST* `/api/v1/todos`

Request Body (JSON):

```json
{
    "title": "My new todo item",
    "completed": false
}
```

> `title` is required.

completed is optional (defaults to false).

Response (201 Cr# eat
ed):

**```json
**{
*   * `"id": "genera`ted-uuid",
    "title": "My new todo item",
    "completed": false
}
```json
```

Curl Command:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"title": "Learn Go Gin", "completed": false}' http://localhost:8080/api/v1/todos
```

#### Get All TODOs

Retriev```es all TODO items.

**Endpoint**: *GET* `/api/v1/todos`

```bash
Response (200 OK):
```
```json
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
```

Curl Command:

```bash
curl http://localhost:8080/api/v1/todos
```

#### Get TODO by ID

Retrieves a single TODO item by its ID.

**Endpoint**: *GET* `/api/v1/todos/:id`

Response (200 OK):

```json
{
    "id": "uuid-of-todo",
    "title": "Specific Todo Item",
    "completed": false
}
```

Response (404 Not Found):

```json
{
    "error": "Todo not found"
}
```

Curl Command (replace `:id` with an actual TODO ID):

```bash
curl http://localhost:8080/api/v1/todos/:id
```

#### Update TODO

Updates an existing TODO item by its ID.

**Endpoint**: *PUT* `/api/v1/todos/:id`

Request Body (JSON):

```json
{
    "title": "Updated todo title",
    "completed": true
}
```

Provide only the fields you want to update. title is required if present.

Response (200 OK):

```json
{
    "id": "uuid-of-todo",
    "title": "Updated todo title",
    "completed": true
}
```

Response (404 Not Found):

```json
{
    "error": "Todo not found"
}
```

Curl Command (replace [ID] with an actual TODO ID):

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"title": "Updated title", "completed": true}' http://localhost:8080/api/v1/todos/[ID]
```

#### Delete TODO

Deletes a TODO item by its ID.

**Endpoint**: *DELETE* `/api/v1/todos/:id`

Response (204 No Content): Successful deletion.

Response (404 Not Found):

```json
{
    "error": "Todo not found"
}
```

Curl Command (replace [ID] with an actual TODO ID):

```json
curl -X DELETE http://localhost:8080/api/v1/todos/[ID]
```

### Testing

Backend (API)

To run the unit tests for the API handlers:

```bash
cd api
go test ./...
```
