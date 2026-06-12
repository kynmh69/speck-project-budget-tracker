package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/project-budget-tracker/backend/internal/config"
	"github.com/your-org/project-budget-tracker/backend/internal/database"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	seedUserEmail    = "demo@example.com"
	seedUserPassword = "password123"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := seed(database.GetDB()); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
}

func seed(db *gorm.DB) error {
	// 冪等性: デモユーザーが既に存在する場合はスキップ
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", seedUserEmail).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Seed data already exists (user %s found), skipping", seedUserEmail)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		user, err := seedUser(tx)
		if err != nil {
			return err
		}

		members, err := seedMembers(tx)
		if err != nil {
			return err
		}

		if err := seedProjects(tx, user, members); err != nil {
			return err
		}

		log.Println("Seed data created successfully")
		log.Printf("Login: %s / %s", seedUserEmail, seedUserPassword)
		return nil
	})
}

func seedUser(tx *gorm.DB) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(seedUserPassword), 12)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        seedUserEmail,
		PasswordHash: string(hash),
		Name:         "デモユーザー",
		Role:         "admin",
	}
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	log.Printf("Created user: %s", user.Email)
	return user, nil
}

func seedMembers(tx *gorm.DB) ([]models.Member, error) {
	members := []models.Member{
		{Name: "佐藤 太郎", Email: "sato@example.com", Role: ptr("PM"), HourlyRate: 8000, Department: ptr("開発部")},
		{Name: "鈴木 花子", Email: "suzuki@example.com", Role: ptr("エンジニア"), HourlyRate: 6000, Department: ptr("開発部")},
		{Name: "田中 一郎", Email: "tanaka@example.com", Role: ptr("デザイナー"), HourlyRate: 5000, Department: ptr("デザイン部")},
	}
	if err := tx.Create(&members).Error; err != nil {
		return nil, err
	}
	log.Printf("Created %d members", len(members))
	return members, nil
}

type taskSeed struct {
	name         string
	status       string
	plannedHours float64
	actualHours  float64
	assignee     *models.Member
}

type projectSeed struct {
	name         string
	description  string
	status       string
	budgetAmount float64
	revenue      float64
	startDate    time.Time
	endDate      time.Time
	tasks        []taskSeed
}

func seedProjects(tx *gorm.DB, user *models.User, members []models.Member) error {
	today := time.Now().Truncate(24 * time.Hour)
	pm, eng, des := &members[0], &members[1], &members[2]

	projects := []projectSeed{
		{
			name:         "ECサイトリニューアル",
			description:  "既存ECサイトのフルリニューアル。フロントエンド刷新とバックエンドAPI化。",
			status:       "in_progress",
			budgetAmount: 5000000,
			revenue:      6000000,
			startDate:    today.AddDate(0, -2, 0),
			endDate:      today.AddDate(0, 2, 0),
			tasks: []taskSeed{
				{name: "要件定義", status: "completed", plannedHours: 40, actualHours: 36, assignee: pm},
				{name: "画面デザイン", status: "completed", plannedHours: 60, actualHours: 72, assignee: des},
				{name: "API実装", status: "in_progress", plannedHours: 120, actualHours: 80, assignee: eng},
				{name: "フロントエンド実装", status: "todo", plannedHours: 100, actualHours: 0, assignee: eng},
			},
		},
		{
			name:         "モバイルアプリ開発",
			description:  "新規モバイルアプリのMVP開発。",
			status:       "planning",
			budgetAmount: 3000000,
			revenue:      0,
			startDate:    today.AddDate(0, 1, 0),
			endDate:      today.AddDate(0, 4, 0),
			tasks: []taskSeed{
				{name: "技術選定", status: "in_progress", plannedHours: 24, actualHours: 8, assignee: eng},
				{name: "プロトタイプ作成", status: "todo", plannedHours: 80, actualHours: 0, assignee: des},
			},
		},
		{
			name:         "社内システム保守",
			description:  "社内勤怠システムの保守・改修対応。",
			status:       "completed",
			budgetAmount: 2000000,
			revenue:      2400000,
			startDate:    today.AddDate(0, -6, 0),
			endDate:      today.AddDate(0, -1, 0),
			tasks: []taskSeed{
				{name: "障害対応", status: "completed", plannedHours: 40, actualHours: 55, assignee: eng},
				{name: "機能改修", status: "completed", plannedHours: 80, actualHours: 76, assignee: eng},
				{name: "運用ドキュメント整備", status: "completed", plannedHours: 20, actualHours: 18, assignee: pm},
			},
		},
	}

	for _, p := range projects {
		if err := seedProject(tx, user, members, p); err != nil {
			return err
		}
	}
	return nil
}

func seedProject(tx *gorm.DB, user *models.User, members []models.Member, p projectSeed) error {
	project := &models.Project{
		UserID:       user.ID,
		Name:         p.name,
		Description:  ptr(p.description),
		Status:       p.status,
		BudgetAmount: ptr(p.budgetAmount),
		StartDate:    ptr(p.startDate),
		EndDate:      ptr(p.endDate),
	}
	if err := tx.Create(project).Error; err != nil {
		return err
	}

	for i := range members {
		assignment := &models.ProjectMember{
			ProjectID:          project.ID,
			MemberID:           members[i].ID,
			JoinedAt:           p.startDate,
			AllocationRate:     1.0,
			HourlyRateSnapshot: ptr(members[i].HourlyRate),
		}
		if err := tx.Create(assignment).Error; err != nil {
			return err
		}
	}

	totalCost, err := seedTasks(tx, user, project, p)
	if err != nil {
		return err
	}

	budget := &models.Budget{
		ProjectID: project.ID,
		Revenue:   p.revenue,
		TotalCost: totalCost,
		Currency:  "JPY",
	}
	budget.CalculateProfit()
	if err := tx.Create(budget).Error; err != nil {
		return err
	}

	log.Printf("Created project: %s (%d tasks, cost=%.0f)", project.Name, len(p.tasks), totalCost)
	return nil
}

// seedTasks creates tasks and time entries; returns the total labor cost
// accumulated from the generated time entries.
func seedTasks(tx *gorm.DB, user *models.User, project *models.Project, p projectSeed) (float64, error) {
	var totalCost float64

	for _, t := range p.tasks {
		var assignedTo *uuid.UUID
		if t.assignee != nil {
			assignedTo = ptr(t.assignee.ID)
		}
		task := &models.Task{
			ProjectID:    project.ID,
			AssignedTo:   assignedTo,
			Name:         t.name,
			PlannedHours: t.plannedHours,
			ActualHours:  t.actualHours,
			Status:       t.status,
			StartDate:    p.StartDateOrNil(),
			EndDate:      p.EndDateOrNil(),
		}
		if err := tx.Create(task).Error; err != nil {
			return 0, err
		}

		cost, err := seedTimeEntries(tx, user, task, t)
		if err != nil {
			return 0, err
		}
		totalCost += cost
	}
	return totalCost, nil
}

// seedTimeEntries splits a task's actual hours into 8-hour daily entries
// working backwards from yesterday.
func seedTimeEntries(tx *gorm.DB, user *models.User, task *models.Task, t taskSeed) (float64, error) {
	if t.assignee == nil || t.actualHours <= 0 {
		return 0, nil
	}

	var totalCost float64
	remaining := t.actualHours
	workDate := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	for remaining > 0 {
		hours := remaining
		if hours > 8 {
			hours = 8
		}
		entry := &models.TimeEntry{
			TaskID:             task.ID,
			MemberID:           t.assignee.ID,
			UserID:             user.ID,
			WorkDate:           workDate,
			Hours:              hours,
			HourlyRateSnapshot: ptr(t.assignee.HourlyRate),
			Comment:            ptr(task.Name + "の作業"),
		}
		if err := tx.Create(entry).Error; err != nil {
			return 0, err
		}
		totalCost += entry.Cost()
		remaining -= hours
		workDate = workDate.AddDate(0, 0, -1)
	}
	return totalCost, nil
}

func (p projectSeed) StartDateOrNil() *time.Time {
	if p.startDate.IsZero() {
		return nil
	}
	return ptr(p.startDate)
}

func (p projectSeed) EndDateOrNil() *time.Time {
	if p.endDate.IsZero() {
		return nil
	}
	return ptr(p.endDate)
}

func ptr[T any](v T) *T {
	return &v
}
