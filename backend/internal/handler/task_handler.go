package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateTask handles POST /api/v1/projects/:projectId/tasks
func (h *TaskHandler) CreateTask(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	var req dto.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_FAILED", err.Error(), nil))
	}

	task, err := h.taskService.CreateTask(projectID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(task))
}

// GetTask handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetTask(c echo.Context) error {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid task ID", nil))
	}

	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(task))
}

// ListTasks handles GET /api/v1/projects/:projectId/tasks
func (h *TaskHandler) ListTasks(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	// Parse pagination params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	status := c.QueryParam("status")

	tasks, err := h.taskService.ListTasksByProject(projectID, page, perPage, status)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(tasks))
}

// UpdateTask handles PUT /api/v1/tasks/:id
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid task ID", nil))
	}

	var req dto.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	task, err := h.taskService.UpdateTask(taskID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(task))
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid task ID", nil))
	}

	if err := h.taskService.DeleteTask(taskID); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{"message": "Task deleted successfully"}))
}

// GetProjectSummary handles GET /api/v1/projects/:id/summary
func (h *TaskHandler) GetProjectSummary(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	summary, err := h.taskService.GetProjectSummary(projectID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(summary))
}

// handleError converts AppError to HTTP response
func handleError(c echo.Context, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
	}
	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "An internal error occurred", nil))
}
