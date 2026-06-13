package service

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
)

// Trend periods supported by GetTrends
const (
	TrendPeriodDaily   = "daily"
	TrendPeriodWeekly  = "weekly"
	TrendPeriodMonthly = "monthly"
)

// AnalyticsService handles business logic for chart/analytics data
type AnalyticsService struct {
	db            *gorm.DB
	timeEntryRepo *repository.TimeEntryRepository
}

// NewAnalyticsService creates a new AnalyticsService
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{
		db:            db,
		timeEntryRepo: repository.NewTimeEntryRepository(db),
	}
}

func (s *AnalyticsService) verifyProject(projectID uuid.UUID) error {
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound("Project")
		}
		return apperrors.ErrDatabaseError(err)
	}
	return nil
}

// GetPlanActual returns planned vs actual hours per task
func (s *AnalyticsService) GetPlanActual(projectID uuid.UUID) (*dto.PlanActualResponse, error) {
	if err := s.verifyProject(projectID); err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := s.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&tasks).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	response := &dto.PlanActualResponse{Tasks: make([]dto.PlanActualTaskData, len(tasks))}
	for i, task := range tasks {
		response.Tasks[i] = dto.PlanActualTaskData{
			TaskName:     task.Name,
			PlannedHours: task.PlannedHours,
			ActualHours:  task.ActualHours,
			Variance:     task.VarianceHours(),
		}
	}
	return response, nil
}

// GetBudgetAnalytics returns revenue, cost breakdown by member, and profit
func (s *AnalyticsService) GetBudgetAnalytics(projectID uuid.UUID) (*dto.BudgetAnalyticsResponse, error) {
	if err := s.verifyProject(projectID); err != nil {
		return nil, err
	}

	var revenue float64
	var budget models.Budget
	if err := s.db.First(&budget, "project_id = ?", projectID).Error; err == nil {
		revenue = budget.Revenue
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrDatabaseError(err)
	}

	summary, err := s.timeEntryRepo.GetSummaryByProject(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	memberSummaries, err := s.timeEntryRepo.GetSummaryByMember(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	breakdown := make([]dto.MemberCostData, len(memberSummaries))
	for i, ms := range memberSummaries {
		breakdown[i] = dto.MemberCostData{MemberName: ms.MemberName, Cost: ms.Cost}
	}

	profit := revenue - summary.TotalCost
	profitRate := 0.0
	if revenue > 0 {
		profitRate = (profit / revenue) * 100
	}

	return &dto.BudgetAnalyticsResponse{
		Revenue:       revenue,
		TotalCost:     summary.TotalCost,
		CostBreakdown: breakdown,
		Profit:        profit,
		ProfitRate:    profitRate,
	}, nil
}

// GetTrends returns time-series data bucketed by the given period
// (daily / weekly / monthly). Actual hours and cost come from time
// entries; planned hours are attributed to the bucket containing the
// task's start date.
func (s *AnalyticsService) GetTrends(projectID uuid.UUID, period string) (*dto.TrendResponse, error) {
	switch period {
	case "":
		period = TrendPeriodMonthly
	case TrendPeriodDaily, TrendPeriodWeekly, TrendPeriodMonthly:
	default:
		return nil, apperrors.ErrValidationFailed("period must be one of: daily, weekly, monthly")
	}

	if err := s.verifyProject(projectID); err != nil {
		return nil, err
	}

	entries, err := s.timeEntryRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	var tasks []models.Task
	if err := s.db.Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	buckets := map[string]*dto.TrendDataPoint{}
	bucketOf := func(t time.Time) string { return bucketKey(t, period) }

	for _, entry := range entries {
		key := bucketOf(entry.WorkDate)
		point, ok := buckets[key]
		if !ok {
			point = &dto.TrendDataPoint{Date: key}
			buckets[key] = point
		}
		point.ActualHours += entry.Hours
		point.Cost += entry.Cost()
	}

	for _, task := range tasks {
		if task.StartDate == nil {
			continue
		}
		key := bucketOf(*task.StartDate)
		point, ok := buckets[key]
		if !ok {
			point = &dto.TrendDataPoint{Date: key}
			buckets[key] = point
		}
		point.PlannedHours += task.PlannedHours
	}

	response := &dto.TrendResponse{Period: period, DataPoints: make([]dto.TrendDataPoint, 0, len(buckets))}
	for _, point := range buckets {
		response.DataPoints = append(response.DataPoints, *point)
	}
	sort.Slice(response.DataPoints, func(i, j int) bool {
		return response.DataPoints[i].Date < response.DataPoints[j].Date
	})
	return response, nil
}

// bucketKey maps a date to its bucket label for the given period
func bucketKey(t time.Time, period string) string {
	switch period {
	case TrendPeriodDaily:
		return t.Format("2006-01-02")
	case TrendPeriodWeekly:
		// 週の開始日（月曜）の日付をバケットとして使用
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		return t.AddDate(0, 0, -(weekday - 1)).Format("2006-01-02")
	default: // monthly
		return t.Format("2006-01")
	}
}

// GetTaskDistribution returns each task's share of the project's actual hours
func (s *AnalyticsService) GetTaskDistribution(projectID uuid.UUID) (*dto.TaskDistributionResponse, error) {
	if err := s.verifyProject(projectID); err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := s.db.Where("project_id = ?", projectID).Order("actual_hours DESC").Find(&tasks).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	var totalActual float64
	for _, task := range tasks {
		totalActual += task.ActualHours
	}

	response := &dto.TaskDistributionResponse{Tasks: make([]dto.TaskDistributionItem, len(tasks))}
	for i, task := range tasks {
		percentage := 0.0
		if totalActual > 0 {
			percentage = (task.ActualHours / totalActual) * 100
		}
		response.Tasks[i] = dto.TaskDistributionItem{
			TaskName:    task.Name,
			ActualHours: task.ActualHours,
			Percentage:  percentage,
		}
	}
	return response, nil
}

// GetProjectsComparison returns KPIs for all of the user's projects
func (s *AnalyticsService) GetProjectsComparison(userIDStr string) (*dto.ProjectsComparisonResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	var projects []models.Project
	if err := s.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&projects).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	response := &dto.ProjectsComparisonResponse{
		Projects: make([]dto.ProjectComparisonItem, len(projects)),
	}

	for i, project := range projects {
		item := dto.ProjectComparisonItem{
			ProjectID:   project.ID.String(),
			ProjectName: project.Name,
		}

		var hours struct {
			PlannedHours float64
			ActualHours  float64
		}
		if err := s.db.Model(&models.Task{}).
			Where("project_id = ?", project.ID).
			Select("COALESCE(SUM(planned_hours), 0) AS planned_hours, COALESCE(SUM(actual_hours), 0) AS actual_hours").
			Scan(&hours).Error; err != nil {
			return nil, apperrors.ErrDatabaseError(err)
		}
		item.PlannedHours = hours.PlannedHours
		item.ActualHours = hours.ActualHours

		summary, err := s.timeEntryRepo.GetSummaryByProject(project.ID)
		if err != nil {
			return nil, apperrors.ErrDatabaseError(err)
		}
		item.TotalCost = summary.TotalCost

		var budget models.Budget
		if err := s.db.First(&budget, "project_id = ?", project.ID).Error; err == nil {
			item.Revenue = budget.Revenue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrDatabaseError(err)
		}

		item.Profit = item.Revenue - item.TotalCost
		if item.Revenue > 0 {
			item.ProfitRate = (item.Profit / item.Revenue) * 100
		}

		response.Projects[i] = item
	}
	return response, nil
}
