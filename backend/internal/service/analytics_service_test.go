package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

func setupAnalyticsTestDB(t *testing.T) *gorm.DB {
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
	return db
}

type analyticsFixture struct {
	db      *gorm.DB
	svc     *service.AnalyticsService
	user    *models.User
	project *models.Project
	member  *models.Member
}

func setupAnalyticsFixture(t *testing.T) *analyticsFixture {
	db := setupAnalyticsTestDB(t)

	user := &models.User{ID: uuid.New(), Email: "test@example.com", PasswordHash: "hash", Name: "Test", Role: "member"}
	require.NoError(t, db.Create(user).Error)

	project := &models.Project{ID: uuid.New(), UserID: user.ID, Name: "分析対象プロジェクト", Status: "in_progress"}
	require.NoError(t, db.Create(project).Error)

	member := &models.Member{ID: uuid.New(), Name: "メンバーA", Email: "a@example.com", HourlyRate: 5000}
	require.NoError(t, db.Create(member).Error)

	return &analyticsFixture{db: db, svc: service.NewAnalyticsService(db), user: user, project: project, member: member}
}

func (f *analyticsFixture) createTask(t *testing.T, name string, planned, actual float64, start *time.Time) *models.Task {
	task := &models.Task{
		ID: uuid.New(), ProjectID: f.project.ID, Name: name,
		PlannedHours: planned, ActualHours: actual, Status: "in_progress", StartDate: start,
	}
	require.NoError(t, f.db.Create(task).Error)
	return task
}

func (f *analyticsFixture) createEntry(t *testing.T, task *models.Task, workDate string, hours, rate float64) {
	d, err := time.Parse("2006-01-02", workDate)
	require.NoError(t, err)
	entry := &models.TimeEntry{
		ID: uuid.New(), TaskID: task.ID, MemberID: f.member.ID, UserID: f.user.ID,
		WorkDate: d, Hours: hours, HourlyRateSnapshot: &rate,
	}
	require.NoError(t, f.db.Create(entry).Error)
}

func TestAnalyticsService_GetPlanActual(t *testing.T) {
	f := setupAnalyticsFixture(t)
	f.createTask(t, "設計", 40, 36, nil)
	f.createTask(t, "実装", 80, 100, nil)

	result, err := f.svc.GetPlanActual(f.project.ID)
	require.NoError(t, err)
	require.Len(t, result.Tasks, 2)
	assert.Equal(t, "設計", result.Tasks[0].TaskName)
	assert.Equal(t, -4.0, result.Tasks[0].Variance)
	assert.Equal(t, 20.0, result.Tasks[1].Variance)

	_, err = f.svc.GetPlanActual(uuid.New())
	require.Error(t, err)
}

func TestAnalyticsService_GetBudgetAnalytics(t *testing.T) {
	f := setupAnalyticsFixture(t)
	task := f.createTask(t, "実装", 80, 0, nil)
	f.createEntry(t, task, "2026-05-01", 10, 5000) // cost 50,000
	f.createEntry(t, task, "2026-05-02", 10, 5000) // cost 50,000

	budget := &models.Budget{ID: uuid.New(), ProjectID: f.project.ID, Revenue: 500000, Currency: "JPY"}
	require.NoError(t, f.db.Create(budget).Error)

	result, err := f.svc.GetBudgetAnalytics(f.project.ID)
	require.NoError(t, err)
	assert.Equal(t, 500000.0, result.Revenue)
	assert.Equal(t, 100000.0, result.TotalCost)
	assert.Equal(t, 400000.0, result.Profit)
	assert.Equal(t, 80.0, result.ProfitRate)
	require.Len(t, result.CostBreakdown, 1)
	assert.Equal(t, "メンバーA", result.CostBreakdown[0].MemberName)
	assert.Equal(t, 100000.0, result.CostBreakdown[0].Cost)
}

func TestAnalyticsService_GetTrends(t *testing.T) {
	f := setupAnalyticsFixture(t)
	start := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)
	task := f.createTask(t, "実装", 100, 0, &start)
	f.createEntry(t, task, "2026-04-15", 8, 5000)
	f.createEntry(t, task, "2026-04-16", 4, 5000)
	f.createEntry(t, task, "2026-05-01", 6, 5000)

	t.Run("monthly: 月単位で集計される", func(t *testing.T) {
		result, err := f.svc.GetTrends(f.project.ID, "")
		require.NoError(t, err)
		assert.Equal(t, "monthly", result.Period)
		require.Len(t, result.DataPoints, 2)

		april := result.DataPoints[0]
		assert.Equal(t, "2026-04", april.Date)
		assert.Equal(t, 12.0, april.ActualHours)
		assert.Equal(t, 60000.0, april.Cost)
		assert.Equal(t, 100.0, april.PlannedHours) // タスク開始日のバケットに計上

		may := result.DataPoints[1]
		assert.Equal(t, "2026-05", may.Date)
		assert.Equal(t, 6.0, may.ActualHours)
	})

	t.Run("daily: 日単位で集計される", func(t *testing.T) {
		result, err := f.svc.GetTrends(f.project.ID, "daily")
		require.NoError(t, err)
		require.Len(t, result.DataPoints, 4) // 工数3日 + タスク開始日
		assert.Equal(t, "2026-04-10", result.DataPoints[0].Date)
		assert.Equal(t, 100.0, result.DataPoints[0].PlannedHours)
	})

	t.Run("weekly: 週の開始日（月曜）に丸められる", func(t *testing.T) {
		result, err := f.svc.GetTrends(f.project.ID, "weekly")
		require.NoError(t, err)
		// 2026-04-15(水)・16(木) は 2026-04-13(月) 週に入る
		var found bool
		for _, p := range result.DataPoints {
			if p.Date == "2026-04-13" {
				assert.Equal(t, 12.0, p.ActualHours)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("不正なperiodはエラー", func(t *testing.T) {
		_, err := f.svc.GetTrends(f.project.ID, "yearly")
		require.Error(t, err)
	})
}

func TestAnalyticsService_GetTaskDistribution(t *testing.T) {
	f := setupAnalyticsFixture(t)
	f.createTask(t, "実装", 0, 75, nil)
	f.createTask(t, "設計", 0, 25, nil)
	f.createTask(t, "未着手", 0, 0, nil)

	result, err := f.svc.GetTaskDistribution(f.project.ID)
	require.NoError(t, err)
	require.Len(t, result.Tasks, 3)
	assert.Equal(t, "実装", result.Tasks[0].TaskName)
	assert.Equal(t, 75.0, result.Tasks[0].Percentage)
	assert.Equal(t, 25.0, result.Tasks[1].Percentage)
	assert.Equal(t, 0.0, result.Tasks[2].Percentage)
}

func TestAnalyticsService_GetProjectsComparison(t *testing.T) {
	f := setupAnalyticsFixture(t)
	task := f.createTask(t, "実装", 50, 40, nil)
	f.createEntry(t, task, "2026-05-01", 20, 5000) // cost 100,000

	budget := &models.Budget{ID: uuid.New(), ProjectID: f.project.ID, Revenue: 300000, Currency: "JPY"}
	require.NoError(t, f.db.Create(budget).Error)

	// 収支未登録の2つ目のプロジェクト
	project2 := &models.Project{ID: uuid.New(), UserID: f.user.ID, Name: "2つ目", Status: "planning"}
	require.NoError(t, f.db.Create(project2).Error)

	result, err := f.svc.GetProjectsComparison(f.user.ID.String())
	require.NoError(t, err)
	require.Len(t, result.Projects, 2)

	first := result.Projects[0]
	assert.Equal(t, "分析対象プロジェクト", first.ProjectName)
	assert.Equal(t, 50.0, first.PlannedHours)
	assert.Equal(t, 40.0, first.ActualHours)
	assert.Equal(t, 300000.0, first.Revenue)
	assert.Equal(t, 100000.0, first.TotalCost)
	assert.Equal(t, 200000.0, first.Profit)

	second := result.Projects[1]
	assert.Equal(t, 0.0, second.Revenue)
	assert.Equal(t, 0.0, second.Profit)

	_, err = f.svc.GetProjectsComparison("not-a-uuid")
	require.Error(t, err)
}
