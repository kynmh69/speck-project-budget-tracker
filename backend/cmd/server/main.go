package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/your-org/project-budget-tracker/backend/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"message": "Project Budget Tracker API is running",
		})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")
	v1.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "API v1",
			"version": "1.0.0",
		})
	})

	// Start server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := e.Start(cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
