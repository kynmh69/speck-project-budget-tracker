package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	"github.com/your-org/project-budget-tracker/backend/internal/handler"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

func setupProjectTestServer(t *testing.T) (*echo.Echo, *gorm.DB, *models.User) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.Project{}, &models.Task{}, &models.Member{})
	require.NoError(t, err)

	// Create test user
	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
		Role:         "member",
	}
	require.NoError(t, db.Create(user).Error)

	// Setup Echo server
	e := echo.New()

	// Initialize repository, service and handler
	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	// Register routes with user context middleware
	api := e.Group("/api/v1")
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", user.ID.String())
			return next(c)
		}
	})
	api.POST("/projects", projectHandler.CreateProject)
	api.GET("/projects", projectHandler.ListProjects)
	api.GET("/projects/:id", projectHandler.GetProject)
	api.PUT("/projects/:id", projectHandler.UpdateProject)
	api.DELETE("/projects/:id", projectHandler.DeleteProject)

	return e, db, user
}

func TestProjectAPI_CreateProject(t *testing.T) {
	e, _, _ := setupProjectTestServer(t)

	t.Run("正常系: プロジェクトを作成できる", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "新しいプロジェクト",
			"description": "プロジェクトの説明",
			"status":      "planning",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("正常系: 全フィールド指定でプロジェクトを作成できる", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":          "フルスペックプロジェクト",
			"description":   "詳細な説明文",
			"status":        "in_progress",
			"start_date":    "2024-01-01",
			"end_date":      "2024-12-31",
			"budget_amount": 1000000,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("異常系: 名前なしでエラー", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"description": "プロジェクトの説明",
			"status":      "planning",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestProjectAPI_GetProject(t *testing.T) {
	e, db, user := setupProjectTestServer(t)

	// Create a project
	project := &models.Project{
		ID:          1,
		OwnerID:     uint(user.ID[0])<<24 | uint(user.ID[1])<<16 | uint(user.ID[2])<<8 | uint(user.ID[3]),
		Name:        "テストプロジェクト",
		Description: "プロジェクトの説明",
		Status:      "planning",
	}
	require.NoError(t, db.Create(project).Error)

	t.Run("正常系: プロジェクトを取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/projects/%d", project.ID), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("異常系: 存在しないプロジェクトIDでエラー", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/9999", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("異常系: 無効なIDフォーマット", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/invalid", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestProjectAPI_ListProjects(t *testing.T) {
	e, db, user := setupProjectTestServer(t)

	// Create multiple projects
	for i := 0; i < 15; i++ {
		status := "planning"
		if i%3 == 0 {
			status = "in_progress"
		}
		project := &models.Project{
			OwnerID:     uint(user.ID[0])<<24 | uint(user.ID[1])<<16 | uint(user.ID[2])<<8 | uint(user.ID[3]),
			Name:        fmt.Sprintf("プロジェクト %d", i+1),
			Description: fmt.Sprintf("プロジェクト %d の説明", i+1),
			Status:      status,
		}
		require.NoError(t, db.Create(project).Error)
	}

	t.Run("正常系: プロジェクト一覧を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("正常系: ページネーションが機能する", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?page=1&per_page=5", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("正常系: ステータスでフィルタできる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?status=in_progress", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("正常系: キーワードで検索できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?keyword=プロジェクト 1", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestProjectAPI_UpdateProject(t *testing.T) {
	e, db, user := setupProjectTestServer(t)

	// Create a project
	project := &models.Project{
		ID:          1,
		OwnerID:     uint(user.ID[0])<<24 | uint(user.ID[1])<<16 | uint(user.ID[2])<<8 | uint(user.ID[3]),
		Name:        "更新前プロジェクト",
		Description: "更新前の説明",
		Status:      "planning",
	}
	require.NoError(t, db.Create(project).Error)

	t.Run("正常系: プロジェクトを更新できる", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "更新後プロジェクト",
			"description": "更新後の説明",
			"status":      "in_progress",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/projects/%d", project.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("異常系: 存在しないプロジェクトの更新でエラー", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "更新後プロジェクト",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/projects/9999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestProjectAPI_DeleteProject(t *testing.T) {
	e, db, user := setupProjectTestServer(t)

	// Create a project
	project := &models.Project{
		ID:          1,
		OwnerID:     uint(user.ID[0])<<24 | uint(user.ID[1])<<16 | uint(user.ID[2])<<8 | uint(user.ID[3]),
		Name:        "削除対象プロジェクト",
		Status:      "planning",
	}
	require.NoError(t, db.Create(project).Error)

	t.Run("正常系: プロジェクトを削除できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/projects/%d", project.ID), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify deletion
		var deleted models.Project
		result := db.First(&deleted, project.ID)
		assert.Error(t, result.Error)
	})

	t.Run("異常系: 存在しないプロジェクトの削除でエラー", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/projects/9999", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("異常系: 無効なIDフォーマット", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/projects/invalid", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
