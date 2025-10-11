package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
	customvalidator "github.com/your-org/project-budget-tracker/backend/internal/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.Response
// @Failure 409 {object} dto.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	// Validate request
	if err := customvalidator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", "Validation failed", err.Error()))
	}

	// Register user
	user, token, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "Internal server error", nil))
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse(response))
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", nil))
	}

	// Validate request
	if err := customvalidator.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", "Validation failed", err.Error()))
	}

	// Login user
	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "Internal server error", nil))
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

// Me godoc
// @Summary Get current user
// @Tags auth
// @Produce json
// @Success 200 {object} dto.UserInfo
// @Failure 401 {object} dto.Response
// @Security BearerAuth
// @Router /auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "User not authenticated", nil))
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.ErrorResponse(appErr.Code, appErr.Message, nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("INTERNAL_ERROR", "Internal server error", nil))
	}

	userInfo := dto.UserInfo{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse(userInfo))
}
