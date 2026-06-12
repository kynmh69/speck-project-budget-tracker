package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
)

const recentProjectsLimit = 5

// DashboardService handles business logic for the dashboard summary
type DashboardService struct {
	db *gorm.DB
}

// NewDashboardService creates a new DashboardService
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

// GetDashboard aggregates project counts, budget totals, and recent projects
// for the given user
func (s *DashboardService) GetDashboard(userIDStr string) (*dto.DashboardResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	response := &dto.DashboardResponse{}

	counts := []struct {
		status string
		dest   *int64
	}{
		{"", &response.TotalProjects},
		{"in_progress", &response.ActiveProjects},
		{"completed", &response.CompletedProjects},
	}
	for _, c := range counts {
		query := s.db.Model(&models.Project{}).Where("user_id = ?", userID)
		if c.status != "" {
			query = query.Where("status = ?", c.status)
		}
		if err := query.Count(c.dest).Error; err != nil {
			return nil, apperrors.ErrDatabaseError(err)
		}
	}

	if err := s.aggregateBudgets(userID, response); err != nil {
		return nil, err
	}

	if err := s.loadRecentProjects(userID, response); err != nil {
		return nil, err
	}

	return response, nil
}

// aggregateBudgets sums revenue/profit over the stored budgets of the user's
// projects. Profit rates are averaged only over projects with revenue, so
// projects without revenue do not drag the average down.
func (s *DashboardService) aggregateBudgets(userID uuid.UUID, response *dto.DashboardResponse) error {
	var agg struct {
		TotalRevenue      float64
		TotalProfit       float64
		AverageProfitRate float64
	}

	err := s.db.Model(&models.Budget{}).
		Joins("JOIN projects ON projects.id = budgets.project_id").
		Where("projects.user_id = ? AND projects.deleted_at IS NULL", userID).
		Select(`COALESCE(SUM(budgets.revenue), 0) AS total_revenue,
			COALESCE(SUM(budgets.profit), 0) AS total_profit,
			COALESCE(AVG(CASE WHEN budgets.revenue > 0 THEN budgets.profit_rate END), 0) AS average_profit_rate`).
		Scan(&agg).Error
	if err != nil {
		return apperrors.ErrDatabaseError(err)
	}

	response.TotalRevenue = agg.TotalRevenue
	response.TotalProfit = agg.TotalProfit
	response.AverageProfitRate = agg.AverageProfitRate
	return nil
}

func (s *DashboardService) loadRecentProjects(userID uuid.UUID, response *dto.DashboardResponse) error {
	var projects []models.Project
	err := s.db.
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(recentProjectsLimit).
		Find(&projects).Error
	if err != nil {
		return apperrors.ErrDatabaseError(err)
	}

	response.RecentProjects = make([]dto.ProjectResponse, len(projects))
	for i := range projects {
		response.RecentProjects[i] = *toDashboardProjectResponse(&projects[i])
	}
	return nil
}

// toDashboardProjectResponse converts a Project model to ProjectResponse DTO
func toDashboardProjectResponse(project *models.Project) *dto.ProjectResponse {
	response := &dto.ProjectResponse{
		ID:        project.ID,
		UserID:    project.UserID,
		Name:      project.Name,
		Status:    project.Status,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}

	if project.Description != nil && *project.Description != "" {
		response.Description = project.Description
	}

	if project.BudgetAmount != nil {
		response.BudgetAmount = project.BudgetAmount
	}

	if project.StartDate != nil {
		formatted := project.StartDate.Format("2006-01-02")
		response.StartDate = &formatted
	}

	if project.EndDate != nil {
		formatted := project.EndDate.Format("2006-01-02")
		response.EndDate = &formatted
	}

	return response
}
