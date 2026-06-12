package handler_test

import (
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

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	"github.com/your-org/project-budget-tracker/backend/internal/handler"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

func setupAnalyticsHandlerTest(t *testing.T) (*handler.AnalyticsHandler, *gorm.DB, *models.Project) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// AutoMigrate cannot be used here: the models declare the PostgreSQL-only
	// default uuid_generate_v4(), which sqlite rejects
	ddl := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY, email TEXT NOT NULL, password_hash TEXT NOT NULL,
			name TEXT NOT NULL, role TEXT NOT NULL DEFAULT 'member',
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE projects (
			id TEXT PRIMARY KEY, user_id TEXT NOT NULL, name TEXT NOT NULL,
			description TEXT, status TEXT NOT NULL DEFAULT 'planning',
			budget_amount REAL, start_date DATE, end_date DATE,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE tasks (
			id TEXT PRIMARY KEY, project_id TEXT NOT NULL, assigned_to TEXT,
			name TEXT NOT NULL, description TEXT,
			planned_hours REAL DEFAULT 0, actual_hours REAL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'todo', start_date DATE, end_date DATE,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE members (
			id TEXT PRIMARY KEY, user_id TEXT, name TEXT NOT NULL, email TEXT NOT NULL,
			role TEXT, hourly_rate REAL DEFAULT 0, department TEXT,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE time_entries (
			id TEXT PRIMARY KEY, task_id TEXT NOT NULL, member_id TEXT NOT NULL,
			user_id TEXT NOT NULL, work_date DATE NOT NULL, hours REAL NOT NULL,
			hourly_rate_snapshot REAL, comment TEXT,
			created_at DATETIME, updated_at DATETIME)`,
		`CREATE TABLE budgets (
			id TEXT PRIMARY KEY, project_id TEXT NOT NULL UNIQUE,
			revenue REAL DEFAULT 0, total_cost REAL DEFAULT 0,
			profit REAL DEFAULT 0, profit_rate REAL DEFAULT 0,
			currency TEXT NOT NULL DEFAULT 'JPY',
			created_at DATETIME, updated_at DATETIME)`,
	}
	for _, stmt := range ddl {
		require.NoError(t, db.Exec(stmt).Error)
	}

	user := &models.User{ID: uuid.New(), Email: "t@example.com", PasswordHash: "h", Name: "T", Role: "member"}
	require.NoError(t, db.Create(user).Error)
	project := &models.Project{ID: uuid.New(), UserID: user.ID, Name: "P", Status: "in_progress"}
	require.NoError(t, db.Create(project).Error)
	task := &models.Task{ID: uuid.New(), ProjectID: project.ID, Name: "実装", PlannedHours: 10, ActualHours: 12, Status: "in_progress"}
	require.NoError(t, db.Create(task).Error)

	return handler.NewAnalyticsHandler(service.NewAnalyticsService(db)), db, project
}

func analyticsRequest(e *echo.Echo, path, paramValue string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if paramValue != "" {
		c.SetParamNames("id")
		c.SetParamValues(paramValue)
	}
	return c, rec
}

func TestAnalyticsHandler_GetPlanActual(t *testing.T) {
	e := echo.New()
	h, _, project := setupAnalyticsHandlerTest(t)

	t.Run("正常系: タスク別予実が返る", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/projects/x/analytics/plan-actual", project.ID.String())
		require.NoError(t, h.GetPlanActual(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Success bool                   `json:"success"`
			Data    dto.PlanActualResponse `json:"data"`
		}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp.Data.Tasks, 1)
		assert.Equal(t, 2.0, resp.Data.Tasks[0].Variance)
	})

	t.Run("異常系: 不正なIDは400", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/projects/x/analytics/plan-actual", "not-a-uuid")
		require.NoError(t, h.GetPlanActual(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しないプロジェクトは404", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/projects/x/analytics/plan-actual", uuid.NewString())
		require.NoError(t, h.GetPlanActual(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAnalyticsHandler_GetTrends(t *testing.T) {
	e := echo.New()
	h, _, project := setupAnalyticsHandlerTest(t)

	t.Run("正常系: デフォルトはmonthly", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/projects/x/analytics/trends", project.ID.String())
		require.NoError(t, h.GetTrends(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"period":"monthly"`)
	})

	t.Run("異常系: 不正なperiodは400", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/projects/x/analytics/trends?period=yearly", project.ID.String())
		require.NoError(t, h.GetTrends(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAnalyticsHandler_GetBudgetAnalyticsAndDistribution(t *testing.T) {
	e := echo.New()
	h, _, project := setupAnalyticsHandlerTest(t)

	c, rec := analyticsRequest(e, "/api/v1/projects/x/analytics/budget", project.ID.String())
	require.NoError(t, h.GetBudgetAnalytics(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	c, rec = analyticsRequest(e, "/api/v1/projects/x/analytics/task-distribution", project.ID.String())
	require.NoError(t, h.GetTaskDistribution(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"percentage":100`)
}

func TestAnalyticsHandler_GetProjectsComparison(t *testing.T) {
	e := echo.New()
	h, db, project := setupAnalyticsHandlerTest(t)

	var owner models.Project
	require.NoError(t, db.First(&owner, "id = ?", project.ID).Error)

	t.Run("正常系: ユーザーのプロジェクト比較が返る", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/analytics/projects-comparison", "")
		c.Set("user_id", owner.UserID.String())
		require.NoError(t, h.GetProjectsComparison(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), project.Name)
	})

	t.Run("異常系: 認証コンテキストなしは401", func(t *testing.T) {
		c, rec := analyticsRequest(e, "/api/v1/analytics/projects-comparison", "")
		require.NoError(t, h.GetProjectsComparison(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
