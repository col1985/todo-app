package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors" // Import the CORS middleware
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"                           // For generating unique IDs
	ginprometheus "github.com/zsais/go-gin-prometheus" // Import gin-prometheus for metrics

	"my-todo-app-api/handlers" // Custom handlers package
	"my-todo-app-api/models"   // Custom models package
	"my-todo-app-api/storage"  // Custom storage package
)

func main() {
	// Set Gin to release mode for production, or debug for development
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS middleware
	// This allows requests from your React frontend (http://localhost:3000)
	// to access your Go Gin API.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // Allowed request headers
		ExposeHeaders:    []string{"Content-Length"}, // Headers to expose to the client
		AllowCredentials: true, // Allow cookies, authorization headers, etc.
		MaxAge:           12 * time.Hour, // How long the preflight request can be cached
	}))

		// Initialize Prometheus middleware
	// This will expose metrics at /metrics endpoint
	p := ginprometheus.NewPrometheus("my_todo_api")
	// Exclude health check endpoint from metrics if desired
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		if c.Request.URL.Path == "/health" {
			return "" // Exclude /health from metrics
		}
		return c.Request.URL.Path
	}
	p.Use(router) // Apply the middleware to the router

	// Initialize in-memory storage
	// This map will hold our Todo items, keyed by their ID
	// In a real application, you would use a database here (e.g., PostgreSQL, MongoDB)
	storage.Todos = make(map[string]models.Todo)

	router.Static("/docs", "./static") // Serves files from ./static at /docs path

// Redirect root to /docs for convenience
	router.GET("/", func(c *gin.Context) {
    	c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	// Define API routes
	// Group routes under /api/v1 for versioning
	v1 := router.Group("/api/v1")
	{
		// POST /api/v1/todos - Create a new TODO item
		v1.POST("/todos", handlers.CreateTodo)

		// GET /api/v1/todos - Get all TODO items
		v1.GET("/todos", handlers.GetTodos)

		// GET /api/v1/todos/:id - Get a single TODO item by ID
		v1.GET("/todos/:id", handlers.GetTodoByID)

		// PUT /api/v1/todos/:id - Update an existing TODO item by ID
		v1.PUT("/todos/:id", handlers.UpdateTodo)

		// DELETE /api/v1/todos/:id - Delete a TODO item by ID
		v1.DELETE("/todos/:id", handlers.DeleteTodo)
	}

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start the server
	// You can configure the port here, e.g., ":8080"
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// init function to pre-populate some data for demonstration
func init() {
	// Ensure the storage map is initialized before use in init
	if storage.Todos == nil {
		storage.Todos = make(map[string]models.Todo)
	}

	// Add some dummy data
	id1 := uuid.New().String()
	storage.Todos[id1] = models.Todo{
		ID:        id1,
		Title:     "Learn Go Gin",
		Completed: false,
	}

	id2 := uuid.New().String()
	storage.Todos[id2] = models.Todo{
		ID:        id2,
		Title:     "Build a REST API",
		Completed: true,
	}

	id3 := uuid.New().String()
	storage.Todos[id3] = models.Todo{
		ID:        id3,
		Title:     "Write Unit Tests",
		Completed: false,
	}
}
