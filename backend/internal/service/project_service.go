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

// ProjectService handles business logic for projects
type ProjectService struct {
	projectRepo *repository.ProjectRepository
	db          *gorm.DB
}

// NewProjectService creates a new ProjectService
func NewProjectService(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// NewProjectServiceWithDB creates a new ProjectService with DB (for tests)
func NewProjectServiceWithDB(db *gorm.DB) *ProjectService {
	projectRepo := repository.NewProjectRepository(db)
	return &ProjectService{
		projectRepo: projectRepo,
		db:          db,
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(userIDStr string, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	project := &models.Project{
		ID:     uuid.New(),
		UserID: userID,
		Name:   req.Name,
		Status: req.Status,
	}

	if project.Status == "" {
		project.Status = "planning"
	}

	if req.Description != nil {
		project.Description = req.Description
	}

	if req.BudgetAmount != nil {
		project.BudgetAmount = req.BudgetAmount
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			project.StartDate = &startDate
		}
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			project.EndDate = &endDate
		}
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return s.toProjectResponse(project), nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(projectIDStr, userIDStr string) (*dto.ProjectDetailResponse, error) {
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	project, err := s.projectRepo.GetByIDWithDetails(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("project")
		}
		return nil, err
	}

	// Check ownership
	isOwner, err := s.projectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrForbidden()
	}

	// Get project stats
	stats, err := s.projectRepo.GetProjectStats(projectID)
	if err != nil {
		return nil, err
	}

	response := &dto.ProjectDetailResponse{
		ProjectResponse: *s.toProjectResponse(project),
	}

	if stats != nil {
		response.Stats = &dto.ProjectStatsResponse{
			TotalTasks:        stats.TotalTasks,
			CompletedTasks:    stats.CompletedTasks,
			TotalPlannedHours: stats.TotalPlannedHours,
			TotalActualHours:  stats.TotalActualHours,
			CompletionRate:    stats.CompletionRate,
		}
	}

	return response, nil
}

// ListProjects retrieves projects with pagination and filtering
func (s *ProjectService) ListProjects(userIDStr string, params dto.ProjectListParams) (*dto.ProjectListResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PerPage <= 0 {
		params.PerPage = 10
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}

	repoParams := repository.ProjectListParams{
		UserID:  userID,
		Page:    params.Page,
		PerPage: params.PerPage,
		Status:  params.Status,
		Search:  params.Search,
		Sort:    params.Sort,
		Order:   params.Order,
	}

	projects, total, err := s.projectRepo.List(repoParams)
	if err != nil {
		return nil, err
	}

	projectResponses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = *s.toProjectResponse(&project)
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage > 0 {
		totalPages++
	}

	return &dto.ProjectListResponse{
		Projects: projectResponses,
		Pagination: dto.Pagination{
			Page:       params.Page,
			PerPage:    params.PerPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(projectIDStr, userIDStr string, req dto.UpdateProjectRequest) (*dto.ProjectDetailResponse, error) {
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("project")
		}
		return nil, err
	}

	// Check ownership
	isOwner, err := s.projectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrForbidden()
	}

	// Update fields
	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.BudgetAmount != nil {
		project.BudgetAmount = req.BudgetAmount
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			project.StartDate = &startDate
		}
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			project.EndDate = &endDate
		}
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return s.GetProject(projectIDStr, userIDStr)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(projectIDStr, userIDStr string) error {
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return apperrors.ErrInvalidInput(err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return apperrors.ErrInvalidInput(err)
	}

	// Check if project exists
	_, err = s.projectRepo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound("project")
		}
		return err
	}

	// Check ownership
	isOwner, err := s.projectRepo.IsOwner(projectID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperrors.ErrForbidden()
	}

	return s.projectRepo.Delete(projectID)
}

// toProjectResponse converts a Project model to ProjectResponse DTO
func (s *ProjectService) toProjectResponse(project *models.Project) *dto.ProjectResponse {
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
