package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
)

// BudgetService handles business logic for budget management
type BudgetService struct {
	db              *gorm.DB
	timeEntryRepo   *repository.TimeEntryRepository
	memberRepo      *repository.MemberRepository
}

// NewBudgetService creates a new BudgetService
func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{
		db:              db,
		timeEntryRepo:   repository.NewTimeEntryRepository(db),
		memberRepo:      repository.NewMemberRepository(db),
	}
}

// GetBudget retrieves or creates a budget for a project
func (s *BudgetService) GetBudget(projectID uuid.UUID) (*dto.BudgetResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Get or create budget
	var budget models.Budget
	if err := s.db.FirstOrCreate(&budget, models.Budget{ProjectID: projectID}).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Calculate current cost from time entries
	summary, err := s.timeEntryRepo.GetSummaryByProject(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Update total cost and recalculate profit
	budget.TotalCost = summary.TotalCost
	budget.CalculateProfit()

	if err := s.db.Save(&budget).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toBudgetResponse(&budget), nil
}

// UpdateRevenue updates the revenue for a project
func (s *BudgetService) UpdateRevenue(projectID uuid.UUID, req *dto.UpdateRevenueRequest) (*dto.BudgetResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Get or create budget
	var budget models.Budget
	if err := s.db.FirstOrCreate(&budget, models.Budget{ProjectID: projectID}).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Update revenue
	budget.Revenue = req.Revenue
	if req.Currency != nil {
		budget.Currency = *req.Currency
	}

	// Recalculate profit
	budget.CalculateProfit()

	if err := s.db.Save(&budget).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toBudgetResponse(&budget), nil
}

// GetBudgetSummary retrieves a comprehensive budget summary for a project
func (s *BudgetService) GetBudgetSummary(projectID uuid.UUID) (*dto.BudgetSummaryResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Get budget
	budget, err := s.GetBudget(projectID)
	if err != nil {
		return nil, err
	}

	// Get time entry summary
	summary, err := s.timeEntryRepo.GetSummaryByProject(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Get cost breakdown by member
	memberSummaries, err := s.timeEntryRepo.GetSummaryByMember(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Calculate average rate
	averageRate := 0.0
	if summary.TotalHours > 0 {
		averageRate = summary.TotalCost / summary.TotalHours
	}

	// Convert member summaries to response
	memberCosts := make([]dto.MemberCostResponse, len(memberSummaries))
	for i, ms := range memberSummaries {
		percentage := 0.0
		if summary.TotalCost > 0 {
			percentage = (ms.Cost / summary.TotalCost) * 100
		}
		memberCosts[i] = dto.MemberCostResponse{
			MemberID:   ms.MemberID,
			MemberName: ms.MemberName,
			Hours:      ms.Hours,
			HourlyRate: ms.HourlyRate,
			Cost:       ms.Cost,
			Percentage: percentage,
		}
	}

	// Create warning message if deficit
	var warningMessage *string
	if budget.IsDeficit {
		msg := "警告: このプロジェクトは赤字です。収益を増やすかコストを削減する必要があります。"
		warningMessage = &msg
	}

	return &dto.BudgetSummaryResponse{
		ProjectID:   projectID,
		ProjectName: project.Name,
		Budget:      *budget,
		CostBreakdown: dto.CostBreakdownResponse{
			LaborCost:   summary.TotalCost,
			TotalHours:  summary.TotalHours,
			AverageRate: averageRate,
		},
		MemberCosts:    memberCosts,
		WarningMessage: warningMessage,
	}, nil
}

// CreateTimeEntry creates a new time entry
func (s *BudgetService) CreateTimeEntry(userID uuid.UUID, req *dto.CreateTimeEntryRequest) (*dto.TimeEntryResponse, error) {
	// Verify task exists
	var task models.Task
	if err := s.db.First(&task, "id = ?", req.TaskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Task")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Verify member exists and get hourly rate
	member, err := s.memberRepo.GetByID(req.MemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Member")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Parse work date
	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	// Create time entry with hourly rate snapshot
	hourlyRate := member.HourlyRate
	timeEntry := &models.TimeEntry{
		TaskID:             req.TaskID,
		MemberID:           req.MemberID,
		UserID:             userID,
		WorkDate:           workDate,
		Hours:              req.Hours,
		HourlyRateSnapshot: &hourlyRate,
		Comment:            req.Comment,
	}

	if err := s.timeEntryRepo.Create(timeEntry); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Update task actual hours
	task.ActualHours += req.Hours
	if err := s.db.Save(&task).Error; err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Reload time entry with relations
	entry, err := s.timeEntryRepo.GetByID(timeEntry.ID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toTimeEntryResponse(entry), nil
}

// GetTimeEntry retrieves a time entry by ID
func (s *BudgetService) GetTimeEntry(id uuid.UUID) (*dto.TimeEntryResponse, error) {
	entry, err := s.timeEntryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("TimeEntry")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toTimeEntryResponse(entry), nil
}

// ListTimeEntries retrieves time entries with filtering
func (s *BudgetService) ListTimeEntries(params repository.TimeEntryListParams) (*dto.TimeEntryListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 || params.PerPage > 100 {
		params.PerPage = 20
	}

	entries, total, err := s.timeEntryRepo.List(params)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Convert to response
	entryResponses := make([]dto.TimeEntryResponse, len(entries))
	for i, entry := range entries {
		entryResponses[i] = *s.toTimeEntryResponse(&entry)
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage > 0 {
		totalPages++
	}

	// Calculate summary
	var totalHours, totalCost float64
	for _, entry := range entries {
		totalHours += entry.Hours
		totalCost += entry.Cost()
	}

	return &dto.TimeEntryListResponse{
		TimeEntries: entryResponses,
		Pagination: dto.Pagination{
			Page:       params.Page,
			PerPage:    params.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
		Summary: &dto.TimeEntrySummary{
			TotalHours: totalHours,
			TotalCost:  totalCost,
		},
	}, nil
}

// UpdateTimeEntry updates a time entry
func (s *BudgetService) UpdateTimeEntry(id uuid.UUID, req *dto.UpdateTimeEntryRequest) (*dto.TimeEntryResponse, error) {
	entry, err := s.timeEntryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("TimeEntry")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Track hours change for task update
	oldHours := entry.Hours

	// Update fields
	if req.WorkDate != nil {
		workDate, err := time.Parse("2006-01-02", *req.WorkDate)
		if err != nil {
			return nil, apperrors.ErrInvalidInput(err)
		}
		entry.WorkDate = workDate
	}
	if req.Hours != nil {
		entry.Hours = *req.Hours
	}
	if req.Comment != nil {
		entry.Comment = req.Comment
	}

	if err := s.timeEntryRepo.Update(entry); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Update task actual hours if hours changed
	if req.Hours != nil && *req.Hours != oldHours {
		var task models.Task
		if err := s.db.First(&task, "id = ?", entry.TaskID).Error; err == nil {
			task.ActualHours = task.ActualHours - oldHours + *req.Hours
			s.db.Save(&task)
		}
	}

	return s.toTimeEntryResponse(entry), nil
}

// DeleteTimeEntry deletes a time entry
func (s *BudgetService) DeleteTimeEntry(id uuid.UUID) error {
	entry, err := s.timeEntryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound("TimeEntry")
		}
		return apperrors.ErrDatabaseError(err)
	}

	// Update task actual hours
	var task models.Task
	if err := s.db.First(&task, "id = ?", entry.TaskID).Error; err == nil {
		task.ActualHours -= entry.Hours
		if task.ActualHours < 0 {
			task.ActualHours = 0
		}
		s.db.Save(&task)
	}

	if err := s.timeEntryRepo.Delete(id); err != nil {
		return apperrors.ErrDatabaseError(err)
	}

	return nil
}

// toBudgetResponse converts a Budget model to BudgetResponse DTO
func (s *BudgetService) toBudgetResponse(budget *models.Budget) *dto.BudgetResponse {
	return &dto.BudgetResponse{
		ID:         budget.ID,
		ProjectID:  budget.ProjectID,
		Revenue:    budget.Revenue,
		TotalCost:  budget.TotalCost,
		Profit:     budget.Profit,
		ProfitRate: budget.ProfitRate,
		Currency:   budget.Currency,
		IsDeficit:  budget.Profit < 0,
	}
}

// toTimeEntryResponse converts a TimeEntry model to TimeEntryResponse DTO
func (s *BudgetService) toTimeEntryResponse(entry *models.TimeEntry) *dto.TimeEntryResponse {
	response := &dto.TimeEntryResponse{
		ID:                 entry.ID,
		TaskID:             entry.TaskID,
		MemberID:           entry.MemberID,
		UserID:             entry.UserID,
		WorkDate:           entry.WorkDate.Format("2006-01-02"),
		Hours:              entry.Hours,
		HourlyRateSnapshot: entry.HourlyRateSnapshot,
		Cost:               entry.Cost(),
		Comment:            entry.Comment,
		CreatedAt:          entry.CreatedAt,
		UpdatedAt:          entry.UpdatedAt,
	}

	if entry.Member.ID != uuid.Nil {
		response.Member = &dto.MemberBriefResponse{
			ID:   entry.Member.ID,
			Name: entry.Member.Name,
		}
	}

	return response
}
