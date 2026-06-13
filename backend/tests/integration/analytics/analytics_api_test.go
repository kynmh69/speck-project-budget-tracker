// Package analytics_test contains integration tests for the analytics API.
//
// NOTE: tests/integration 直下のパッケージはコンパイルエラーのため、
// 独立してビルド・実行できるようサブパッケージとして分離している。
package analytics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
)

// setupAnalyticsSchema はSQLite互換のテーブルスキーマを作成
// (モデルのuuid_generate_v4()デフォルトがPostgreSQL専用のためAutoMigrate不可)
func setupAnalyticsSchema(t *testing.T, db *gorm.DB) {
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
}

type fixture struct {
	e       *echo.Echo
	db      *gorm.DB
	user    *models.User
	project *models.Project
}

// setupAnalyticsServer は本番のルーティング構成（認証ミドルウェアが
// user_idをstringで格納する）を再現したテストサーバーを構築する
func setupAnalyticsServer(t *testing.T) *fixture {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	setupAnalyticsSchema(t, db)

	user := &models.User{ID: uuid.New(), Email: "test@example.com", PasswordHash: "hash", Name: "Test User", Role: "member"}
	require.NoError(t, db.Create(user).Error)

	project := &models.Project{ID: uuid.New(), UserID: user.ID, Name: "分析プロジェクト", Status: "in_progress"}
	require.NoError(t, db.Create(project).Error)

	member := &models.Member{ID: uuid.New(), Name: "メンバーA", Email: "a@example.com", HourlyRate: 5000}
	require.NoError(t, db.Create(member).Error)

	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	task1 := &models.Task{ID: uuid.New(), ProjectID: project.ID, Name: "設計", PlannedHours: 40, ActualHours: 30, Status: "completed", StartDate: &start}
	task2 := &models.Task{ID: uuid.New(), ProjectID: project.ID, Name: "実装", PlannedHours: 60, ActualHours: 90, Status: "in_progress", StartDate: &start}
	require.NoError(t, db.Create(task1).Error)
	require.NoError(t, db.Create(task2).Error)

	rate := 5000.0
	entries := []struct {
		task *models.Task
		date time.Time
		hrs  float64
	}{
		{task1, time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC), 10},
		{task2, time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), 8},
		{task2, time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC), 6},
	}
	for _, e := range entries {
		entry := &models.TimeEntry{
			ID: uuid.New(), TaskID: e.task.ID, MemberID: member.ID, UserID: user.ID,
			WorkDate: e.date, Hours: e.hrs, HourlyRateSnapshot: &rate,
		}
		require.NoError(t, db.Create(entry).Error)
	}

	budget := &models.Budget{ID: uuid.New(), ProjectID: project.ID, Revenue: 1000000, Currency: "JPY"}
	require.NoError(t, db.Create(budget).Error)

	e := echo.New()
	analyticsHandler := handler.NewAnalyticsHandler(service.NewAnalyticsService(db))

	api := e.Group("/api/v1")
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 本番のAuthMiddlewareと同様にstringで格納
			c.Set("user_id", user.ID.String())
			return next(c)
		}
	})
	api.GET("/projects/:id/analytics/plan-actual", analyticsHandler.GetPlanActual)
	api.GET("/projects/:id/analytics/budget", analyticsHandler.GetBudgetAnalytics)
	api.GET("/projects/:id/analytics/trends", analyticsHandler.GetTrends)
	api.GET("/projects/:id/analytics/task-distribution", analyticsHandler.GetTaskDistribution)
	api.GET("/analytics/projects-comparison", analyticsHandler.GetProjectsComparison)

	return &fixture{e: e, db: db, user: user, project: project}
}

func get[T any](t *testing.T, f *fixture, path string) (int, T) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	f.e.ServeHTTP(rec, req)

	var resp struct {
		Success bool `json:"success"`
		Data    T    `json:"data"`
	}
	if rec.Code == http.StatusOK {
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	}
	return rec.Code, resp.Data
}

func TestAnalyticsAPI_PlanActual(t *testing.T) {
	f := setupAnalyticsServer(t)

	code, data := get[dto.PlanActualResponse](t, f, "/api/v1/projects/"+f.project.ID.String()+"/analytics/plan-actual")
	require.Equal(t, http.StatusOK, code)
	require.Len(t, data.Tasks, 2)
	assert.Equal(t, -10.0, data.Tasks[0].Variance) // 設計: 30 - 40
	assert.Equal(t, 30.0, data.Tasks[1].Variance)  // 実装: 90 - 60
}

func TestAnalyticsAPI_Budget(t *testing.T) {
	f := setupAnalyticsServer(t)

	code, data := get[dto.BudgetAnalyticsResponse](t, f, "/api/v1/projects/"+f.project.ID.String()+"/analytics/budget")
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, 1000000.0, data.Revenue)
	assert.Equal(t, 120000.0, data.TotalCost) // 24h × 5000
	assert.Equal(t, 880000.0, data.Profit)
	require.Len(t, data.CostBreakdown, 1)
	assert.Equal(t, "メンバーA", data.CostBreakdown[0].MemberName)
}

func TestAnalyticsAPI_Trends(t *testing.T) {
	f := setupAnalyticsServer(t)

	code, data := get[dto.TrendResponse](t, f, "/api/v1/projects/"+f.project.ID.String()+"/analytics/trends?period=monthly")
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, "monthly", data.Period)
	require.Len(t, data.DataPoints, 2)
	assert.Equal(t, "2026-04", data.DataPoints[0].Date)
	assert.Equal(t, 18.0, data.DataPoints[0].ActualHours)
	assert.Equal(t, 100.0, data.DataPoints[0].PlannedHours) // タスク2件の開始月
	assert.Equal(t, "2026-05", data.DataPoints[1].Date)
	assert.Equal(t, 6.0, data.DataPoints[1].ActualHours)

	code, _ = get[dto.TrendResponse](t, f, "/api/v1/projects/"+f.project.ID.String()+"/analytics/trends?period=yearly")
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestAnalyticsAPI_TaskDistribution(t *testing.T) {
	f := setupAnalyticsServer(t)

	code, data := get[dto.TaskDistributionResponse](t, f, "/api/v1/projects/"+f.project.ID.String()+"/analytics/task-distribution")
	require.Equal(t, http.StatusOK, code)
	require.Len(t, data.Tasks, 2)
	assert.Equal(t, "実装", data.Tasks[0].TaskName) // actual_hours降順
	assert.InDelta(t, 75.0, data.Tasks[0].Percentage, 0.01)
	assert.InDelta(t, 25.0, data.Tasks[1].Percentage, 0.01)
}

func TestAnalyticsAPI_ProjectsComparison(t *testing.T) {
	f := setupAnalyticsServer(t)

	// 他ユーザーのプロジェクトは含まれないことを確認するためのデータ
	other := &models.User{ID: uuid.New(), Email: "other@example.com", PasswordHash: "h", Name: "Other", Role: "member"}
	require.NoError(t, f.db.Create(other).Error)
	otherProject := &models.Project{ID: uuid.New(), UserID: other.ID, Name: "他人のプロジェクト", Status: "planning"}
	require.NoError(t, f.db.Create(otherProject).Error)

	code, data := get[dto.ProjectsComparisonResponse](t, f, "/api/v1/analytics/projects-comparison")
	require.Equal(t, http.StatusOK, code)
	require.Len(t, data.Projects, 1)
	item := data.Projects[0]
	assert.Equal(t, "分析プロジェクト", item.ProjectName)
	assert.Equal(t, 100.0, item.PlannedHours)
	assert.Equal(t, 120.0, item.ActualHours)
	assert.Equal(t, 1000000.0, item.Revenue)
	assert.Equal(t, 120000.0, item.TotalCost)
}

func TestAnalyticsAPI_NotFound(t *testing.T) {
	f := setupAnalyticsServer(t)

	code, _ := get[dto.PlanActualResponse](t, f, "/api/v1/projects/"+uuid.NewString()+"/analytics/plan-actual")
	assert.Equal(t, http.StatusNotFound, code)
}
