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

func setupDashboardTestDB(t *testing.T) *gorm.DB {
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

func createDashboardTestUser(t *testing.T, db *gorm.DB) *models.User {
	user := &models.User{
		ID:           uuid.New(),
		Email:        uuid.NewString() + "@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
		Role:         "member",
	}
	require.NoError(t, db.Create(user).Error)
	return user
}

func createDashboardTestProject(t *testing.T, db *gorm.DB, userID uuid.UUID, name, status string) *models.Project {
	project := &models.Project{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
		Status: status,
	}
	require.NoError(t, db.Create(project).Error)
	return project
}

func createDashboardTestBudget(t *testing.T, db *gorm.DB, projectID uuid.UUID, revenue, totalCost float64) {
	budget := &models.Budget{
		ID:        uuid.New(),
		ProjectID: projectID,
		Revenue:   revenue,
		TotalCost: totalCost,
		Currency:  "JPY",
	}
	budget.CalculateProfit()
	require.NoError(t, db.Create(budget).Error)
}

func TestDashboardService_GetDashboard(t *testing.T) {
	db := setupDashboardTestDB(t)
	svc := service.NewDashboardService(db)
	user := createDashboardTestUser(t, db)

	t.Run("正常系: プロジェクトがない場合はゼロ値を返す", func(t *testing.T) {
		result, err := svc.GetDashboard(user.ID.String())
		require.NoError(t, err)
		assert.Equal(t, int64(0), result.TotalProjects)
		assert.Equal(t, int64(0), result.ActiveProjects)
		assert.Equal(t, int64(0), result.CompletedProjects)
		assert.Equal(t, 0.0, result.TotalRevenue)
		assert.Equal(t, 0.0, result.TotalProfit)
		assert.Equal(t, 0.0, result.AverageProfitRate)
		assert.Empty(t, result.RecentProjects)
	})

	t.Run("正常系: ステータス別件数と収支が集計される", func(t *testing.T) {
		p1 := createDashboardTestProject(t, db, user.ID, "進行中プロジェクト", "in_progress")
		p2 := createDashboardTestProject(t, db, user.ID, "完了プロジェクト", "completed")
		createDashboardTestProject(t, db, user.ID, "計画中プロジェクト", "planning")

		// p1: 売上100万・コスト60万 → 利益40万・利益率40%
		createDashboardTestBudget(t, db, p1.ID, 1000000, 600000)
		// p2: 売上200万・コスト100万 → 利益100万・利益率50%
		createDashboardTestBudget(t, db, p2.ID, 2000000, 1000000)

		result, err := svc.GetDashboard(user.ID.String())
		require.NoError(t, err)
		assert.Equal(t, int64(3), result.TotalProjects)
		assert.Equal(t, int64(1), result.ActiveProjects)
		assert.Equal(t, int64(1), result.CompletedProjects)
		assert.Equal(t, 3000000.0, result.TotalRevenue)
		assert.Equal(t, 1400000.0, result.TotalProfit)
		assert.InDelta(t, 45.0, result.AverageProfitRate, 0.01) // (40+50)/2
		assert.Len(t, result.RecentProjects, 3)
	})

	t.Run("正常系: 売上ゼロのプロジェクトは利益率平均に含まれない", func(t *testing.T) {
		db := setupDashboardTestDB(t)
		svc := service.NewDashboardService(db)
		user := createDashboardTestUser(t, db)

		p1 := createDashboardTestProject(t, db, user.ID, "売上あり", "in_progress")
		p2 := createDashboardTestProject(t, db, user.ID, "売上なし", "planning")
		createDashboardTestBudget(t, db, p1.ID, 1000000, 800000) // 利益率20%
		createDashboardTestBudget(t, db, p2.ID, 0, 100000)       // 売上ゼロ

		result, err := svc.GetDashboard(user.ID.String())
		require.NoError(t, err)
		assert.InDelta(t, 20.0, result.AverageProfitRate, 0.01)
	})

	t.Run("正常系: 他ユーザーのプロジェクトは集計されない", func(t *testing.T) {
		other := createDashboardTestUser(t, db)
		otherProject := createDashboardTestProject(t, db, other.ID, "他人のプロジェクト", "in_progress")
		createDashboardTestBudget(t, db, otherProject.ID, 9999999, 0)

		result, err := svc.GetDashboard(other.ID.String())
		require.NoError(t, err)
		assert.Equal(t, int64(1), result.TotalProjects)
		assert.Equal(t, 9999999.0, result.TotalRevenue)
	})

	t.Run("正常系: 最近のプロジェクトは更新日時降順で最大5件", func(t *testing.T) {
		db := setupDashboardTestDB(t)
		svc := service.NewDashboardService(db)
		user := createDashboardTestUser(t, db)

		base := time.Now().Add(-1 * time.Hour)
		for i := 0; i < 7; i++ {
			p := createDashboardTestProject(t, db, user.ID, "プロジェクト", "planning")
			require.NoError(t, db.Model(p).Update("updated_at", base.Add(time.Duration(i)*time.Minute)).Error)
		}

		result, err := svc.GetDashboard(user.ID.String())
		require.NoError(t, err)
		require.Len(t, result.RecentProjects, 5)
		for i := 1; i < len(result.RecentProjects); i++ {
			assert.True(t, !result.RecentProjects[i].UpdatedAt.After(result.RecentProjects[i-1].UpdatedAt))
		}
	})

	t.Run("異常系: 不正なユーザーIDでエラー", func(t *testing.T) {
		_, err := svc.GetDashboard("not-a-uuid")
		require.Error(t, err)
	})
}
