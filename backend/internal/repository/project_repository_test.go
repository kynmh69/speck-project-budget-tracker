package repository_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
)

func setupProjectRepoTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// AutoMigrate cannot be used here: the models declare the PostgreSQL-only
	// default uuid_generate_v4(), which sqlite rejects
	ddl := []string{
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

func date(s string) *time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return &t
}

func createFilterTestProject(t *testing.T, db *gorm.DB, userID uuid.UUID, name, status string, start, end *time.Time, profitRate *float64) *models.Project {
	project := &models.Project{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Status:    status,
		StartDate: start,
		EndDate:   end,
	}
	require.NoError(t, db.Create(project).Error)

	if profitRate != nil {
		budget := &models.Budget{
			ID:         uuid.New(),
			ProjectID:  project.ID,
			ProfitRate: *profitRate,
			Currency:   "JPY",
		}
		require.NoError(t, db.Create(budget).Error)
	}
	return project
}

func TestProjectRepository_List_ExtendedFilters(t *testing.T) {
	db := setupProjectRepoTestDB(t)
	repo := repository.NewProjectRepository(db)
	userID := uuid.New()
	rate := func(v float64) *float64 { return &v }

	// 2026-01〜2026-03 / 利益率40%
	createFilterTestProject(t, db, userID, "Q1プロジェクト", "completed", date("2026-01-01"), date("2026-03-31"), rate(40))
	// 2026-04〜2026-06 / 利益率10%
	createFilterTestProject(t, db, userID, "Q2プロジェクト", "in_progress", date("2026-04-01"), date("2026-06-30"), rate(10))
	// 期間未設定 / 収支未登録
	createFilterTestProject(t, db, userID, "期間未定プロジェクト", "planning", nil, nil, nil)

	baseParams := func() repository.ProjectListParams {
		return repository.ProjectListParams{UserID: userID, Page: 1, PerPage: 10}
	}

	t.Run("期間フィルタ: 重なるプロジェクトのみ返す", func(t *testing.T) {
		params := baseParams()
		params.PeriodFrom = date("2026-04-01")
		params.PeriodTo = date("2026-12-31")

		projects, total, err := repo.List(params)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total) // Q2 + 期間未定（オープン期間は常に一致）
		names := []string{projects[0].Name, projects[1].Name}
		assert.Contains(t, names, "Q2プロジェクト")
		assert.Contains(t, names, "期間未定プロジェクト")
	})

	t.Run("利益率フィルタ: 下限以上のみ返す", func(t *testing.T) {
		params := baseParams()
		params.MinProfitRate = rate(20)

		projects, total, err := repo.List(params)
		require.NoError(t, err)
		require.Equal(t, int64(1), total)
		assert.Equal(t, "Q1プロジェクト", projects[0].Name)
	})

	t.Run("利益率フィルタ: 収支未登録は0%として扱う", func(t *testing.T) {
		params := baseParams()
		params.MaxProfitRate = rate(5)

		projects, total, err := repo.List(params)
		require.NoError(t, err)
		require.Equal(t, int64(1), total)
		assert.Equal(t, "期間未定プロジェクト", projects[0].Name)
	})

	t.Run("利益率ソート: 降順で返す", func(t *testing.T) {
		params := baseParams()
		params.Sort = "profit_rate"

		projects, _, err := repo.List(params)
		require.NoError(t, err)
		require.Len(t, projects, 3)
		assert.Equal(t, "Q1プロジェクト", projects[0].Name)
		assert.Equal(t, "Q2プロジェクト", projects[1].Name)
		assert.Equal(t, "期間未定プロジェクト", projects[2].Name)
	})

	t.Run("ステータス + 期間の組み合わせ", func(t *testing.T) {
		params := baseParams()
		params.Status = "in_progress"
		params.PeriodFrom = date("2026-01-01")

		projects, total, err := repo.List(params)
		require.NoError(t, err)
		require.Equal(t, int64(1), total)
		assert.Equal(t, "Q2プロジェクト", projects[0].Name)
	})
}
