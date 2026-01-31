package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
)

// TimeEntryRepository handles database operations for time entries
type TimeEntryRepository struct {
	db *gorm.DB
}

// NewTimeEntryRepository creates a new TimeEntryRepository
func NewTimeEntryRepository(db *gorm.DB) *TimeEntryRepository {
	return &TimeEntryRepository{db: db}
}

// Create creates a new time entry
func (r *TimeEntryRepository) Create(entry *models.TimeEntry) error {
	return r.db.Create(entry).Error
}

// GetByID retrieves a time entry by ID
func (r *TimeEntryRepository) GetByID(id uuid.UUID) (*models.TimeEntry, error) {
	var entry models.TimeEntry
	if err := r.db.Preload("Member").Preload("Task").First(&entry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

// TimeEntryListParams represents parameters for listing time entries
type TimeEntryListParams struct {
	ProjectID *uuid.UUID
	TaskID    *uuid.UUID
	MemberID  *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PerPage   int
}

// List retrieves time entries with filtering and pagination
func (r *TimeEntryRepository) List(params TimeEntryListParams) ([]models.TimeEntry, int64, error) {
	var entries []models.TimeEntry
	var total int64

	query := r.db.Model(&models.TimeEntry{})

	// Apply filters
	if params.TaskID != nil {
		query = query.Where("task_id = ?", *params.TaskID)
	}

	if params.MemberID != nil {
		query = query.Where("member_id = ?", *params.MemberID)
	}

	if params.ProjectID != nil {
		query = query.Joins("JOIN tasks ON tasks.id = time_entries.task_id").
			Where("tasks.project_id = ?", *params.ProjectID)
	}

	if params.StartDate != nil {
		query = query.Where("work_date >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("work_date <= ?", *params.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PerPage
	if err := query.
		Preload("Member").
		Preload("Task").
		Order("work_date DESC").
		Offset(offset).
		Limit(params.PerPage).
		Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// Update updates a time entry
func (r *TimeEntryRepository) Update(entry *models.TimeEntry) error {
	return r.db.Save(entry).Error
}

// Delete deletes a time entry
func (r *TimeEntryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.TimeEntry{}, "id = ?", id).Error
}

// GetByProjectID retrieves all time entries for a project
func (r *TimeEntryRepository) GetByProjectID(projectID uuid.UUID) ([]models.TimeEntry, error) {
	var entries []models.TimeEntry
	if err := r.db.
		Joins("JOIN tasks ON tasks.id = time_entries.task_id").
		Where("tasks.project_id = ?", projectID).
		Preload("Member").
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

// GetSummaryByProject calculates total hours and cost for a project
func (r *TimeEntryRepository) GetSummaryByProject(projectID uuid.UUID) (*TimeEntrySummary, error) {
	var summary TimeEntrySummary

	if err := r.db.Model(&models.TimeEntry{}).
		Select(`
			COALESCE(SUM(time_entries.hours), 0) as total_hours,
			COALESCE(SUM(time_entries.hours * COALESCE(time_entries.hourly_rate_snapshot, 0)), 0) as total_cost
		`).
		Joins("JOIN tasks ON tasks.id = time_entries.task_id").
		Where("tasks.project_id = ?", projectID).
		Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

// GetSummaryByMember calculates hours and cost grouped by member for a project
func (r *TimeEntryRepository) GetSummaryByMember(projectID uuid.UUID) ([]MemberCostSummary, error) {
	var summaries []MemberCostSummary

	if err := r.db.Model(&models.TimeEntry{}).
		Select(`
			time_entries.member_id,
			members.name as member_name,
			COALESCE(SUM(time_entries.hours), 0) as hours,
			COALESCE(AVG(time_entries.hourly_rate_snapshot), 0) as hourly_rate,
			COALESCE(SUM(time_entries.hours * COALESCE(time_entries.hourly_rate_snapshot, 0)), 0) as cost
		`).
		Joins("JOIN tasks ON tasks.id = time_entries.task_id").
		Joins("JOIN members ON members.id = time_entries.member_id").
		Where("tasks.project_id = ?", projectID).
		Group("time_entries.member_id, members.name").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

// TimeEntrySummary represents aggregated time entry data
type TimeEntrySummary struct {
	TotalHours float64 `json:"total_hours"`
	TotalCost  float64 `json:"total_cost"`
}

// MemberCostSummary represents cost summary by member
type MemberCostSummary struct {
	MemberID   uuid.UUID `json:"member_id"`
	MemberName string    `json:"member_name"`
	Hours      float64   `json:"hours"`
	HourlyRate float64   `json:"hourly_rate"`
	Cost       float64   `json:"cost"`
}
