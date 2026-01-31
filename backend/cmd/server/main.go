package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/your-org/project-budget-tracker/backend/internal/config"
	"github.com/your-org/project-budget-tracker/backend/internal/database"
	"github.com/your-org/project-budget-tracker/backend/internal/handler"
	custommiddleware "github.com/your-org/project-budget-tracker/backend/internal/middleware"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connected successfully")

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(custommiddleware.CORSConfig())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Project Budget Tracker API is running",
		})
	})

	// Initialize services
	authService := service.NewAuthService(database.GetDB(), cfg.JWTSecret)
	projectRepo := repository.NewProjectRepository(database.GetDB())
	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(database.GetDB())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", authHandler.Me, custommiddleware.AuthMiddleware(authService))

	// Protected routes
	protected := v1.Group("", custommiddleware.AuthMiddleware(authService))

	// Project routes
	protected.POST("/projects", projectHandler.CreateProject)
	protected.GET("/projects", projectHandler.ListProjects)
	protected.GET("/projects/:id", projectHandler.GetProject)
	protected.PUT("/projects/:id", projectHandler.UpdateProject)
	protected.DELETE("/projects/:id", projectHandler.DeleteProject)

	// Task routes
	protected.POST("/projects/:projectId/tasks", taskHandler.CreateTask)
	protected.GET("/projects/:projectId/tasks", taskHandler.ListTasks)
	protected.GET("/projects/:id/summary", taskHandler.GetProjectSummary)
	protected.GET("/tasks/:id", taskHandler.GetTask)
	protected.PUT("/tasks/:id", taskHandler.UpdateTask)
	protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// Start server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := e.Start(cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
