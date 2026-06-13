package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	"github.com/your-org/project-budget-tracker/backend/internal/handler"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
	customvalidator "github.com/your-org/project-budget-tracker/backend/internal/validator"
)

// setupTaskTestDBSchema はSQLite互換のテーブルスキーマを作成
func setupTaskTestDBSchema(t *testing.T, db *gorm.DB) {
	require.NoError(t, db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'member',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error)

	require.NoError(t, db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL DEFAULT 'planning',
			budget_amount REAL,
			start_date DATE,
			end_date DATE,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error)

	require.NoError(t, db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			assigned_to TEXT,
			name TEXT NOT NULL,
			description TEXT,
			planned_hours REAL DEFAULT 0,
			actual_hours REAL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'todo',
			start_date DATE,
			end_date DATE,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error)

	require.NoError(t, db.Exec(`
		CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			role TEXT,
			hourly_rate REAL DEFAULT 0,
			department TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error)
}

func setupTestServer(t *testing.T) (*echo.Echo, *gorm.DB, *models.Project) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// SQLite互換のスキーマを作成
	setupTaskTestDBSchema(t, db)

	// Create test user and project
	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
		Role:         "member",
	}
	require.NoError(t, db.Create(user).Error)

	project := &models.Project{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "Test Project",
		Status: "planning",
	}
	require.NoError(t, db.Create(project).Error)

	// Setup Echo server
	e := echo.New()
	e.Validator = &customvalidator.EchoValidator{}

	// Initialize service and handler
	taskService := service.NewTaskService(db)
	taskHandler := handler.NewTaskHandler(taskService)

	// Register routes
	e.POST("/api/v1/projects/:projectId/tasks", taskHandler.CreateTask)
	e.GET("/api/v1/projects/:projectId/tasks", taskHandler.ListTasks)
	e.GET("/api/v1/projects/:id/summary", taskHandler.GetProjectSummary)
	e.GET("/api/v1/tasks/:id", taskHandler.GetTask)
	e.PUT("/api/v1/tasks/:id", taskHandler.UpdateTask)
	e.DELETE("/api/v1/tasks/:id", taskHandler.DeleteTask)

	return e, db, project
}

func TestTaskAPI_CreateTask(t *testing.T) {
	e, _, project := setupTestServer(t)

	t.Run("正常系: タスクを作成できる", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":          "新しいタスク",
			"description":   "タスクの説明",
			"planned_hours": 8.0,
			"status":        "todo",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/"+project.ID.String()+"/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("異常系: 無効なプロジェクトIDでエラー", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":          "新しいタスク",
			"planned_hours": 8.0,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/invalid-uuid/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestTaskAPI_GetTask(t *testing.T) {
	e, db, project := setupTestServer(t)

	// Create a task
	task := &models.Task{
		ID:           uuid.New(),
		ProjectID:    project.ID,
		Name:         "テストタスク",
		PlannedHours: 10.0,
		ActualHours:  5.0,
		Status:       "in_progress",
	}
	require.NoError(t, db.Create(task).Error)

	t.Run("正常系: タスクを取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+task.ID.String(), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("異常系: 存在しないタスクIDでエラー", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTaskAPI_ListTasks(t *testing.T) {
	e, db, project := setupTestServer(t)

	// Create multiple tasks
	for i := 0; i < 5; i++ {
		task := &models.Task{
			ID:           uuid.New(),
			ProjectID:    project.ID,
			Name:         "タスク" + string(rune('A'+i)),
			PlannedHours: float64(i * 2),
			Status:       "todo",
		}
		require.NoError(t, db.Create(task).Error)
	}

	t.Run("正常系: タスク一覧を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/"+project.ID.String()+"/tasks", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("正常系: ステータスでフィルタリングできる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/"+project.ID.String()+"/tasks?status=todo", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestTaskAPI_UpdateTask(t *testing.T) {
	e, db, project := setupTestServer(t)

	// Create a task
	task := &models.Task{
		ID:           uuid.New(),
		ProjectID:    project.ID,
		Name:         "更新テスト用タスク",
		PlannedHours: 10.0,
		ActualHours:  0.0,
		Status:       "todo",
	}
	require.NoError(t, db.Create(task).Error)

	t.Run("正常系: タスクを更新できる", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"actual_hours": 12.0,
			"status":       "completed",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+task.ID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})
}

func TestTaskAPI_DeleteTask(t *testing.T) {
	e, db, project := setupTestServer(t)

	// Create a task
	task := &models.Task{
		ID:        uuid.New(),
		ProjectID: project.ID,
		Name:      "削除テスト用タスク",
		Status:    "todo",
	}
	require.NoError(t, db.Create(task).Error)

	t.Run("正常系: タスクを削除できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+task.ID.String(), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 削除済みタスクの取得でエラー", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+task.ID.String(), nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTaskAPI_GetProjectSummary(t *testing.T) {
	e, db, project := setupTestServer(t)

	// Create multiple tasks
	tasks := []*models.Task{
		{ID: uuid.New(), ProjectID: project.ID, Name: "完了タスク1", PlannedHours: 10.0, ActualHours: 8.0, Status: "completed"},
		{ID: uuid.New(), ProjectID: project.ID, Name: "完了タスク2", PlannedHours: 20.0, ActualHours: 25.0, Status: "completed"},
		{ID: uuid.New(), ProjectID: project.ID, Name: "進行中タスク", PlannedHours: 15.0, ActualHours: 10.0, Status: "in_progress"},
		{ID: uuid.New(), ProjectID: project.ID, Name: "未着手タスク", PlannedHours: 5.0, ActualHours: 0.0, Status: "todo"},
	}

	for _, task := range tasks {
		require.NoError(t, db.Create(task).Error)
	}

	t.Run("正常系: プロジェクトサマリーを取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/"+project.ID.String()+"/summary", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.Response
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
	})
}
