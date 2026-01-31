package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=200"`
	Description  *string `json:"description,omitempty"`
	Status       string  `json:"status,omitempty" validate:"omitempty,oneof=planning in_progress completed on_hold"`
	BudgetAmount *float64 `json:"budget_amount,omitempty" validate:"omitempty,min=0"`
	StartDate    *string `json:"start_date,omitempty"`
	EndDate      *string `json:"end_date,omitempty"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description  *string  `json:"description,omitempty"`
	Status       *string  `json:"status,omitempty" validate:"omitempty,oneof=planning in_progress completed on_hold"`
	BudgetAmount *float64 `json:"budget_amount,omitempty" validate:"omitempty,min=0"`
	StartDate    *string  `json:"start_date,omitempty"`
	EndDate      *string  `json:"end_date,omitempty"`
}

// ProjectResponse represents a project response
type ProjectResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	Status       string    `json:"status"`
	BudgetAmount *float64  `json:"budget_amount,omitempty"`
	StartDate    *string   `json:"start_date,omitempty"`
	EndDate      *string   `json:"end_date,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProjectDetailResponse represents a detailed project response
type ProjectDetailResponse struct {
	ProjectResponse
	Stats *ProjectStatsResponse `json:"stats,omitempty"`
}

// ProjectStatsResponse represents project statistics
type ProjectStatsResponse struct {
	TotalTasks        int     `json:"total_tasks"`
	CompletedTasks    int     `json:"completed_tasks"`
	TotalPlannedHours float64 `json:"total_planned_hours"`
	TotalActualHours  float64 `json:"total_actual_hours"`
	CompletionRate    float64 `json:"completion_rate"`
}

// ProjectListResponse represents a paginated list of projects
type ProjectListResponse struct {
	Projects   []ProjectResponse `json:"projects"`
	Pagination Pagination        `json:"pagination"`
}

// ProjectListParams represents query parameters for listing projects
type ProjectListParams struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Status  string `query:"status"`
	Search  string `query:"search"`
	Sort    string `query:"sort"`
	Order   string `query:"order"`
}
