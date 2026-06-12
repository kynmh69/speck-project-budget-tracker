// Package dashboard_test contains integration tests for the dashboard API.
//
// NOTE: tests/integration 直下のパッケージはコンパイルエラーのため、
// 独立してビルド・実行できるようサブパッケージとして分離している。
package dashboard_test

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
	"gorm.io/gorm/logger"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	"github.com/your-org/project-budget-tracker/backend/internal/handler"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

// setupDashboardSchema はSQLite互換のテーブルスキーマを作成
// (モデルのuuid_generate_v4()デフォルトがPostgreSQL専用のためAutoMigrate不可)
func setupDashboardSchema(t *testing.T, db *gorm.DB) {
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
}

// setupDashboardServer は本番のルーティング構成（認証ミドルウェアが
// user_idをstringで格納する）を再現したテストサーバーを構築する
func setupDashboardServer(t *testing.T) (*echo.Echo, *gorm.DB, *models.User) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	setupDashboardSchema(t, db)

	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hash",
		Name:         "Test User",
		Role:         "member",
	}
	require.NoError(t, db.Create(user).Error)

	e := echo.New()
	dashboardHandler := handler.NewDashboardHandler(service.NewDashboardService(db))

	api := e.Group("/api/v1")
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 本番のAuthMiddlewareと同様にstringで格納
			c.Set("user_id", user.ID.String())
			return next(c)
		}
	})
	api.GET("/dashboard", dashboardHandler.GetDashboard)

	// 認証なしルート（ミドルウェアなしの比較用）
	e.GET("/api/noauth/dashboard", dashboardHandler.GetDashboard)

	return e, db, user
}

func createProjectWithBudget(t *testing.T, db *gorm.DB, userID uuid.UUID, name, status string, revenue, totalCost float64) {
	project := &models.Project{ID: uuid.New(), UserID: userID, Name: name, Status: status}
	require.NoError(t, db.Create(project).Error)

	if revenue > 0 || totalCost > 0 {
		budget := &models.Budget{ID: uuid.New(), ProjectID: project.ID, Revenue: revenue, TotalCost: totalCost, Currency: "JPY"}
		budget.CalculateProfit()
		require.NoError(t, db.Create(budget).Error)
	}
}

func getDashboard(t *testing.T, e *echo.Echo, path string) (*httptest.ResponseRecorder, dto.DashboardResponse) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var resp struct {
		Success bool                  `json:"success"`
		Data    dto.DashboardResponse `json:"data"`
	}
	if rec.Code == http.StatusOK {
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	}
	return rec, resp.Data
}

func TestDashboardAPI_GetDashboard(t *testing.T) {
	t.Run("正常系: KPIサマリーが返る", func(t *testing.T) {
		e, db, user := setupDashboardServer(t)

		createProjectWithBudget(t, db, user.ID, "進行中A", "in_progress", 1000000, 600000) // 利益40万・利益率40%
		createProjectWithBudget(t, db, user.ID, "進行中B", "in_progress", 0, 0)            // 収支未登録
		createProjectWithBudget(t, db, user.ID, "完了C", "completed", 2000000, 1000000)   // 利益100万・利益率50%
		createProjectWithBudget(t, db, user.ID, "計画D", "planning", 0, 0)

		rec, data := getDashboard(t, e, "/api/v1/dashboard")
		require.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, int64(4), data.TotalProjects)
		assert.Equal(t, int64(2), data.ActiveProjects)
		assert.Equal(t, int64(1), data.CompletedProjects)
		assert.Equal(t, 3000000.0, data.TotalRevenue)
		assert.Equal(t, 1400000.0, data.TotalProfit)
		assert.InDelta(t, 45.0, data.AverageProfitRate, 0.01)
		assert.Len(t, data.RecentProjects, 4)
	})

	t.Run("正常系: 他ユーザーのデータは含まれない", func(t *testing.T) {
		e, db, user := setupDashboardServer(t)

		createProjectWithBudget(t, db, user.ID, "自分のプロジェクト", "in_progress", 500000, 100000)

		otherUser := &models.User{
			ID: uuid.New(), Email: "other@example.com", PasswordHash: "hash", Name: "Other", Role: "member",
		}
		require.NoError(t, db.Create(otherUser).Error)
		createProjectWithBudget(t, db, otherUser.ID, "他人のプロジェクト", "in_progress", 9000000, 0)

		rec, data := getDashboard(t, e, "/api/v1/dashboard")
		require.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, int64(1), data.TotalProjects)
		assert.Equal(t, 500000.0, data.TotalRevenue)
		require.Len(t, data.RecentProjects, 1)
		assert.Equal(t, "自分のプロジェクト", data.RecentProjects[0].Name)
	})

	t.Run("正常系: 最近のプロジェクトは最大5件", func(t *testing.T) {
		e, db, user := setupDashboardServer(t)

		for i := 0; i < 7; i++ {
			createProjectWithBudget(t, db, user.ID, "プロジェクト", "planning", 0, 0)
		}

		rec, data := getDashboard(t, e, "/api/v1/dashboard")
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, int64(7), data.TotalProjects)
		assert.Len(t, data.RecentProjects, 5)
	})

	t.Run("正常系: 論理削除されたプロジェクトは集計対象外", func(t *testing.T) {
		e, db, user := setupDashboardServer(t)

		createProjectWithBudget(t, db, user.ID, "残すプロジェクト", "in_progress", 100000, 0)
		deleted := &models.Project{ID: uuid.New(), UserID: user.ID, Name: "削除予定", Status: "in_progress"}
		require.NoError(t, db.Create(deleted).Error)
		require.NoError(t, db.Delete(deleted).Error)

		rec, data := getDashboard(t, e, "/api/v1/dashboard")
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, int64(1), data.TotalProjects)
		assert.Equal(t, int64(1), data.ActiveProjects)
	})

	t.Run("異常系: 認証コンテキストなしは401", func(t *testing.T) {
		e, _, _ := setupDashboardServer(t)

		rec, _ := getDashboard(t, e, "/api/noauth/dashboard")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
