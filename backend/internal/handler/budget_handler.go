package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

// BudgetHandler handles HTTP requests for budget management
type BudgetHandler struct {
	budgetService *service.BudgetService
}

// NewBudgetHandler creates a new BudgetHandler
func NewBudgetHandler(budgetService *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService}
}

// GetBudget handles GET /api/v1/projects/:id/budget
func (h *BudgetHandler) GetBudget(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	budget, err := h.budgetService.GetBudgetSummary(projectID)
	if err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(budget))
}

// UpdateRevenue handles PUT /api/v1/projects/:id/budget/revenue
func (h *BudgetHandler) UpdateRevenue(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	var req dto.UpdateRevenueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_FAILED", err.Error(), nil))
	}

	budget, err := h.budgetService.UpdateRevenue(projectID, &req)
	if err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(budget))
}

// CreateTimeEntry handles POST /api/v1/time-entries
func (h *BudgetHandler) CreateTimeEntry(c echo.Context) error {
	// Get user ID from context
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "User not authenticated", nil))
	}

	var req dto.CreateTimeEntryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_FAILED", err.Error(), nil))
	}

	entry, err := h.budgetService.CreateTimeEntry(userID, &req)
	if err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(entry))
}

// GetTimeEntry handles GET /api/v1/time-entries/:id
func (h *BudgetHandler) GetTimeEntry(c echo.Context) error {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid time entry ID", nil))
	}

	entry, err := h.budgetService.GetTimeEntry(entryID)
	if err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(entry))
}

// ListTimeEntries handles GET /api/v1/time-entries
func (h *BudgetHandler) ListTimeEntries(c echo.Context) error {
	// Parse pagination params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	params := repository.TimeEntryListParams{
		Page:    page,
		PerPage: perPage,
	}

	// Parse optional filters
	if projectIDStr := c.QueryParam("project_id"); projectIDStr != "" {
		if projectID, err := uuid.Parse(projectIDStr); err == nil {
			params.ProjectID = &projectID
		}
	}

	if taskIDStr := c.QueryParam("task_id"); taskIDStr != "" {
		if taskID, err := uuid.Parse(taskIDStr); err == nil {
			params.TaskID = &taskID
		}
	}

	if memberIDStr := c.QueryParam("member_id"); memberIDStr != "" {
		if memberID, err := uuid.Parse(memberIDStr); err == nil {
			params.MemberID = &memberID
		}
	}

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			params.StartDate = &startDate
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			params.EndDate = &endDate
		}
	}

	entries, err := h.budgetService.ListTimeEntries(params)
	if err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(entries))
}

// UpdateTimeEntry handles PUT /api/v1/time-entries/:id
func (h *BudgetHandler) UpdateTimeEntry(c echo.Context) error {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid time entry ID", nil))
	}

	var req dto.UpdateTimeEntryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	entry, err := h.budgetService.UpdateTimeEntry(entryID, &req)
	if err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(entry))
}

// DeleteTimeEntry handles DELETE /api/v1/time-entries/:id
func (h *BudgetHandler) DeleteTimeEntry(c echo.Context) error {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid time entry ID", nil))
	}

	if err := h.budgetService.DeleteTimeEntry(entryID); err != nil {
		return handleBudgetError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{"message": "Time entry deleted successfully"}))
}

// handleBudgetError converts AppError to HTTP response
func handleBudgetError(c echo.Context, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
	}
	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "An internal error occurred", nil))
}
