package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new ProjectRepository
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByIDWithDetails retrieves a project by ID with related data
func (r *ProjectRepository) GetByIDWithDetails(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.
		Preload("User").
		Preload("Budget").
		First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// ProjectListParams represents parameters for listing projects
type ProjectListParams struct {
	UserID uuid.UUID
	Page   int
	PerPage int
	Status string
	Search string
	Sort   string
	Order  string
}

// List retrieves projects with pagination, filtering, and search
func (r *ProjectRepository) List(params ProjectListParams) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{}).Where("user_id = ?", params.UserID)

	// Apply status filter if provided
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Apply search filter if provided
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Determine sort column
	sortColumn := "created_at"
	switch params.Sort {
	case "name":
		sortColumn = "name"
	case "start_date":
		sortColumn = "start_date"
	case "end_date":
		sortColumn = "end_date"
	case "status":
		sortColumn = "status"
	case "created_at":
		sortColumn = "created_at"
	}

	// Determine sort order
	sortOrder := "DESC"
	if params.Order == "asc" {
		sortOrder = "ASC"
	}

	// Apply pagination and sorting
	offset := (params.Page - 1) * params.PerPage
	if err := query.
		Order(sortColumn + " " + sortOrder).
		Offset(offset).
		Limit(params.PerPage).
		Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// Update updates a project
func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

// Delete soft deletes a project
func (r *ProjectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

// ExistsByID checks if a project exists by ID
func (r *ProjectRepository) ExistsByID(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsOwner checks if a user is the owner of a project
func (r *ProjectRepository) IsOwner(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).
		Where("id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetProjectStats retrieves project statistics
func (r *ProjectRepository) GetProjectStats(projectID uuid.UUID) (*ProjectStats, error) {
	var stats ProjectStats

	if err := r.db.Model(&models.Task{}).
		Select(`
			COUNT(*) as total_tasks,
			COALESCE(SUM(planned_hours), 0) as total_planned_hours,
			COALESCE(SUM(actual_hours), 0) as total_actual_hours,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks
		`).
		Where("project_id = ?", projectID).
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	stats.ProjectID = projectID
	if stats.TotalTasks > 0 {
		stats.CompletionRate = float64(stats.CompletedTasks) / float64(stats.TotalTasks) * 100
	}

	return &stats, nil
}

// ProjectStats represents project statistics
type ProjectStats struct {
	ProjectID          uuid.UUID `json:"project_id"`
	TotalTasks         int       `json:"total_tasks"`
	TotalPlannedHours  float64   `json:"total_planned_hours"`
	TotalActualHours   float64   `json:"total_actual_hours"`
	CompletedTasks     int       `json:"completed_tasks"`
	CompletionRate     float64   `json:"completion_rate"`
}
