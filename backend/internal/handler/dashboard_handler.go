package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

// DashboardHandler handles HTTP requests for the dashboard
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetDashboard handles GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "User not authenticated", nil))
	}

	dashboard, err := h.dashboardService.GetDashboard(userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "An internal error occurred", nil))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(dashboard))
}
