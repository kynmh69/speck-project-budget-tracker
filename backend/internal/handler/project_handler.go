package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
	customvalidator "github.com/your-org/project-budget-tracker/backend/internal/validator"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	projectService *service.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject handles POST /api/v1/projects
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req dto.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "リクエストの形式が正しくありません", nil))
	}

	if err := customvalidator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", "バリデーションエラー", err.Error()))
	}

	project, err := h.projectService.CreateProject(userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(project))
}

// GetProject handles GET /api/v1/projects/:id
func (h *ProjectHandler) GetProject(c echo.Context) error {
	userID := c.Get("user_id").(string)
	projectID := c.Param("id")

	project, err := h.projectService.GetProject(projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(project))
}

// ListProjects handles GET /api/v1/projects
func (h *ProjectHandler) ListProjects(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var params dto.ProjectListParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "パラメータの形式が正しくありません", nil))
	}

	projects, err := h.projectService.ListProjects(userID, params)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(projects))
}

// UpdateProject handles PUT /api/v1/projects/:id
func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	userID := c.Get("user_id").(string)
	projectID := c.Param("id")

	var req dto.UpdateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "リクエストの形式が正しくありません", nil))
	}

	if err := customvalidator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", "バリデーションエラー", err.Error()))
	}

	project, err := h.projectService.UpdateProject(projectID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(project))
}

// DeleteProject handles DELETE /api/v1/projects/:id
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	userID := c.Get("user_id").(string)
	projectID := c.Param("id")

	err := h.projectService.DeleteProject(projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{"message": "プロジェクトを削除しました"}))
}

// handleError handles errors and returns appropriate HTTP responses
func (h *ProjectHandler) handleError(c echo.Context, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
	}

	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "サーバーエラーが発生しました", nil))
}
