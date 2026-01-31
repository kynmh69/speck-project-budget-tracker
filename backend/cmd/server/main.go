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
	memberService := service.NewMemberService(database.GetDB())
	budgetService := service.NewBudgetService(database.GetDB())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)
	memberHandler := handler.NewMemberHandler(memberService)
	budgetHandler := handler.NewBudgetHandler(budgetService)

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

	// Member routes
	protected.POST("/members", memberHandler.CreateMember)
	protected.GET("/members", memberHandler.ListMembers)
	protected.GET("/members/:id", memberHandler.GetMember)
	protected.PUT("/members/:id", memberHandler.UpdateMember)
	protected.DELETE("/members/:id", memberHandler.DeleteMember)

	// Project member routes
	protected.GET("/projects/:id/members", memberHandler.GetProjectMembers)
	protected.POST("/projects/:id/members", memberHandler.AssignMemberToProject)
	protected.DELETE("/projects/:id/members/:memberId", memberHandler.RemoveMemberFromProject)

	// Budget routes
	protected.GET("/projects/:id/budget", budgetHandler.GetBudget)
	protected.PUT("/projects/:id/budget/revenue", budgetHandler.UpdateRevenue)

	// Time entry routes
	protected.POST("/time-entries", budgetHandler.CreateTimeEntry)
	protected.GET("/time-entries", budgetHandler.ListTimeEntries)
	protected.GET("/time-entries/:id", budgetHandler.GetTimeEntry)
	protected.PUT("/time-entries/:id", budgetHandler.UpdateTimeEntry)
	protected.DELETE("/time-entries/:id", budgetHandler.DeleteTimeEntry)

	// Start server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := e.Start(cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
