package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task
func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	if err := r.db.Preload("Assignee").First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByProjectID retrieves tasks by project ID with pagination
func (r *TaskRepository) GetByProjectID(projectID uuid.UUID, page, perPage int, status string) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.Model(&models.Task{}).Where("project_id = ?", projectID)

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * perPage
	if err := query.Preload("Assignee").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// Update updates a task
func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

// Delete soft deletes a task
func (r *TaskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}

// GetProjectSummary calculates the summary of planned and actual hours for a project
func (r *TaskRepository) GetProjectSummary(projectID uuid.UUID) (*TaskSummary, error) {
	var summary TaskSummary

	if err := r.db.Model(&models.Task{}).
		Select(`
			COUNT(*) as total_tasks,
			COALESCE(SUM(planned_hours), 0) as total_planned_hours,
			COALESCE(SUM(actual_hours), 0) as total_actual_hours,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks,
			COUNT(CASE WHEN status = 'todo' THEN 1 END) as todo_tasks,
			COUNT(CASE WHEN status = 'blocked' THEN 1 END) as blocked_tasks
		`).
		Where("project_id = ?", projectID).
		Scan(&summary).Error; err != nil {
		return nil, err
	}

	summary.ProjectID = projectID
	summary.VarianceHours = summary.TotalActualHours - summary.TotalPlannedHours
	if summary.TotalPlannedHours > 0 {
		summary.VariancePercentage = (summary.VarianceHours / summary.TotalPlannedHours) * 100
	}

	return &summary, nil
}

// TaskSummary represents the aggregated task data for a project
type TaskSummary struct {
	ProjectID          uuid.UUID `json:"project_id"`
	TotalTasks         int       `json:"total_tasks"`
	TotalPlannedHours  float64   `json:"total_planned_hours"`
	TotalActualHours   float64   `json:"total_actual_hours"`
	VarianceHours      float64   `json:"variance_hours"`
	VariancePercentage float64   `json:"variance_percentage"`
	CompletedTasks     int       `json:"completed_tasks"`
	InProgressTasks    int       `json:"in_progress_tasks"`
	TodoTasks          int       `json:"todo_tasks"`
	BlockedTasks       int       `json:"blocked_tasks"`
}
