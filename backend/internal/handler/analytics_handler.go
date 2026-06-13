package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

// AnalyticsHandler handles HTTP requests for chart/analytics data
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) projectID(c echo.Context) (uuid.UUID, bool) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, false
	}
	return projectID, true
}

func handleAnalyticsError(c echo.Context, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
	}
	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "An internal error occurred", nil))
}

// GetPlanActual handles GET /api/v1/projects/:id/analytics/plan-actual
func (h *AnalyticsHandler) GetPlanActual(c echo.Context) error {
	projectID, ok := h.projectID(c)
	if !ok {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	result, err := h.analyticsService.GetPlanActual(projectID)
	if err != nil {
		return handleAnalyticsError(c, err)
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetBudgetAnalytics handles GET /api/v1/projects/:id/analytics/budget
func (h *AnalyticsHandler) GetBudgetAnalytics(c echo.Context) error {
	projectID, ok := h.projectID(c)
	if !ok {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	result, err := h.analyticsService.GetBudgetAnalytics(projectID)
	if err != nil {
		return handleAnalyticsError(c, err)
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetTrends handles GET /api/v1/projects/:id/analytics/trends?period=daily|weekly|monthly
func (h *AnalyticsHandler) GetTrends(c echo.Context) error {
	projectID, ok := h.projectID(c)
	if !ok {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	result, err := h.analyticsService.GetTrends(projectID, c.QueryParam("period"))
	if err != nil {
		return handleAnalyticsError(c, err)
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetTaskDistribution handles GET /api/v1/projects/:id/analytics/task-distribution
func (h *AnalyticsHandler) GetTaskDistribution(c echo.Context) error {
	projectID, ok := h.projectID(c)
	if !ok {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	result, err := h.analyticsService.GetTaskDistribution(projectID)
	if err != nil {
		return handleAnalyticsError(c, err)
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetProjectsComparison handles GET /api/v1/analytics/projects-comparison
func (h *AnalyticsHandler) GetProjectsComparison(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "User not authenticated", nil))
	}

	result, err := h.analyticsService.GetProjectsComparison(userID)
	if err != nil {
		return handleAnalyticsError(c, err)
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse(result))
}
