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

// TaskService handles business logic for tasks
type TaskService struct {
	taskRepo *repository.TaskRepository
	db       *gorm.DB
}

// NewTaskService creates a new TaskService
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		taskRepo: repository.NewTaskRepository(db),
		db:       db,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(projectID uuid.UUID, req *dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Parse dates
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, apperrors.ErrInvalidInput(err)
		}
		startDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, apperrors.ErrInvalidInput(err)
		}
		endDate = &t
	}

	// Set default status
	status := "todo"
	if req.Status != "" {
		status = req.Status
	}

	task := &models.Task{
		ProjectID:    projectID,
		AssignedTo:   req.AssignedTo,
		Name:         req.Name,
		Description:  req.Description,
		PlannedHours: req.PlannedHours,
		ActualHours:  0,
		Status:       status,
		StartDate:    startDate,
		EndDate:      endDate,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toTaskResponse(task), nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(id uuid.UUID) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Task")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toTaskResponse(task), nil
}

// ListTasksByProject retrieves tasks for a project with pagination
func (s *TaskService) ListTasksByProject(projectID uuid.UUID, page, perPage int, status string) (*dto.TaskListResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Set default pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	tasks, total, err := s.taskRepo.GetByProjectID(projectID, page, perPage, status)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Convert to response
	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = *s.toTaskResponse(&task)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.TaskListResponse{
		Tasks: taskResponses,
		Pagination: dto.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(id uuid.UUID, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Task")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Update fields
	if req.Name != nil {
		task.Name = *req.Name
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.AssignedTo != nil {
		task.AssignedTo = req.AssignedTo
	}
	if req.PlannedHours != nil {
		task.PlannedHours = *req.PlannedHours
	}
	if req.ActualHours != nil {
		task.ActualHours = *req.ActualHours
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, apperrors.ErrInvalidInput(err)
		}
		task.StartDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, apperrors.ErrInvalidInput(err)
		}
		task.EndDate = &t
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toTaskResponse(task), nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id uuid.UUID) error {
	// Check if task exists
	if _, err := s.taskRepo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound("Task")
		}
		return apperrors.ErrDatabaseError(err)
	}

	if err := s.taskRepo.Delete(id); err != nil {
		return apperrors.ErrDatabaseError(err)
	}

	return nil
}

// GetProjectSummary retrieves the summary of tasks for a project
func (s *TaskService) GetProjectSummary(projectID uuid.UUID) (*dto.ProjectSummaryResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	summary, err := s.taskRepo.GetProjectSummary(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Calculate completion rate
	var completionRate float64
	if summary.TotalTasks > 0 {
		completionRate = float64(summary.CompletedTasks) / float64(summary.TotalTasks) * 100
	}

	return &dto.ProjectSummaryResponse{
		ProjectID:          summary.ProjectID,
		TotalTasks:         summary.TotalTasks,
		TotalPlannedHours:  summary.TotalPlannedHours,
		TotalActualHours:   summary.TotalActualHours,
		VarianceHours:      summary.VarianceHours,
		VariancePercentage: summary.VariancePercentage,
		IsOverBudget:       summary.VarianceHours > 0,
		CompletedTasks:     summary.CompletedTasks,
		InProgressTasks:    summary.InProgressTasks,
		TodoTasks:          summary.TodoTasks,
		BlockedTasks:       summary.BlockedTasks,
		CompletionRate:     completionRate,
	}, nil
}

// toTaskResponse converts a Task model to TaskResponse DTO
func (s *TaskService) toTaskResponse(task *models.Task) *dto.TaskResponse {
	response := &dto.TaskResponse{
		ID:                 task.ID,
		ProjectID:          task.ProjectID,
		AssignedTo:         task.AssignedTo,
		Name:               task.Name,
		Description:        task.Description,
		PlannedHours:       task.PlannedHours,
		ActualHours:        task.ActualHours,
		VarianceHours:      task.VarianceHours(),
		VariancePercentage: task.VariancePercentage(),
		Status:             task.Status,
		CreatedAt:          task.CreatedAt,
		UpdatedAt:          task.UpdatedAt,
	}

	if task.StartDate != nil {
		startDateStr := task.StartDate.Format("2006-01-02")
		response.StartDate = &startDateStr
	}
	if task.EndDate != nil {
		endDateStr := task.EndDate.Format("2006-01-02")
		response.EndDate = &endDateStr
	}

	if task.Assignee != nil {
		response.Assignee = &dto.MemberBriefResponse{
			ID:   task.Assignee.ID,
			Name: task.Assignee.Name,
		}
	}

	return response
}
