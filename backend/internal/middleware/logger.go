package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoggerConfig returns the logger middleware configuration
func LoggerConfig() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	})
}
