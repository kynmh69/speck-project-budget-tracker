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

// MemberHandler handles HTTP requests for members
type MemberHandler struct {
	memberService *service.MemberService
}

// NewMemberHandler creates a new MemberHandler
func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

// CreateMember handles POST /api/v1/members
func (h *MemberHandler) CreateMember(c echo.Context) error {
	var req dto.CreateMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_FAILED", err.Error(), nil))
	}

	member, err := h.memberService.CreateMember(&req)
	if err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(member))
}

// GetMember handles GET /api/v1/members/:id
func (h *MemberHandler) GetMember(c echo.Context) error {
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid member ID", nil))
	}

	member, err := h.memberService.GetMember(memberID)
	if err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(member))
}

// ListMembers handles GET /api/v1/members
func (h *MemberHandler) ListMembers(c echo.Context) error {
	// Parse pagination params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	search := c.QueryParam("search")
	department := c.QueryParam("department")

	members, err := h.memberService.ListMembers(page, perPage, search, department)
	if err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(members))
}

// UpdateMember handles PUT /api/v1/members/:id
func (h *MemberHandler) UpdateMember(c echo.Context) error {
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid member ID", nil))
	}

	var req dto.UpdateMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	member, err := h.memberService.UpdateMember(memberID, &req)
	if err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(member))
}

// DeleteMember handles DELETE /api/v1/members/:id
func (h *MemberHandler) DeleteMember(c echo.Context) error {
	memberID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid member ID", nil))
	}

	if err := h.memberService.DeleteMember(memberID); err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{"message": "Member deleted successfully"}))
}

// GetProjectMembers handles GET /api/v1/projects/:id/members
func (h *MemberHandler) GetProjectMembers(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	members, err := h.memberService.GetProjectMembers(projectID)
	if err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(members))
}

// AssignMemberToProject handles POST /api/v1/projects/:id/members
func (h *MemberHandler) AssignMemberToProject(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	var req dto.AssignMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_FAILED", err.Error(), nil))
	}

	member, err := h.memberService.AssignMemberToProject(projectID, &req)
	if err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(member))
}

// RemoveMemberFromProject handles DELETE /api/v1/projects/:id/members/:memberId
func (h *MemberHandler) RemoveMemberFromProject(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid project ID", nil))
	}

	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid member ID", nil))
	}

	if err := h.memberService.RemoveMemberFromProject(projectID, memberID); err != nil {
		return handleMemberError(c, err)
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{"message": "Member removed from project successfully"}))
}

// handleMemberError converts AppError to HTTP response
func handleMemberError(c echo.Context, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
	}
	return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "An internal error occurred", nil))
}
