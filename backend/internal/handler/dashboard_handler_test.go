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

func setupDashboardHandlerTest(t *testing.T) (*handler.DashboardHandler, *gorm.DB) {
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

	return handler.NewDashboardHandler(service.NewDashboardService(db)), db
}

func newDashboardRequestContext(e *echo.Echo) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestDashboardHandler_GetDashboard(t *testing.T) {
	e := echo.New()

	t.Run("正常系: ダッシュボードサマリーを返す", func(t *testing.T) {
		h, db := setupDashboardHandlerTest(t)

		userID := uuid.New()
		require.NoError(t, db.Create(&models.User{
			ID: userID, Email: "demo@example.com", PasswordHash: "hash", Name: "Demo", Role: "member",
		}).Error)
		project := &models.Project{ID: uuid.New(), UserID: userID, Name: "P1", Status: "in_progress"}
		require.NoError(t, db.Create(project).Error)
		budget := &models.Budget{ID: uuid.New(), ProjectID: project.ID, Revenue: 1000000, TotalCost: 600000, Currency: "JPY"}
		budget.CalculateProfit()
		require.NoError(t, db.Create(budget).Error)

		c, rec := newDashboardRequestContext(e)
		c.Set("user_id", userID.String())

		require.NoError(t, h.GetDashboard(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Success bool                  `json:"success"`
			Data    dto.DashboardResponse `json:"data"`
		}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
		assert.Equal(t, int64(1), resp.Data.TotalProjects)
		assert.Equal(t, int64(1), resp.Data.ActiveProjects)
		assert.Equal(t, 1000000.0, resp.Data.TotalRevenue)
		assert.Equal(t, 400000.0, resp.Data.TotalProfit)
		assert.Len(t, resp.Data.RecentProjects, 1)
	})

	t.Run("異常系: user_id未設定で401", func(t *testing.T) {
		h, _ := setupDashboardHandlerTest(t)

		c, rec := newDashboardRequestContext(e)

		require.NoError(t, h.GetDashboard(c))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("異常系: 不正なuser_idで400", func(t *testing.T) {
		h, _ := setupDashboardHandlerTest(t)

		c, rec := newDashboardRequestContext(e)
		c.Set("user_id", "not-a-uuid")

		require.NoError(t, h.GetDashboard(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
