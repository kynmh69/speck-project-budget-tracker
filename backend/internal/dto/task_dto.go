package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	Name         string     `json:"name" validate:"required,min=1,max=200"`
	Description  *string    `json:"description,omitempty"`
	AssignedTo   *uuid.UUID `json:"assigned_to,omitempty"`
	PlannedHours float64    `json:"planned_hours" validate:"min=0"`
	Status       string     `json:"status,omitempty" validate:"omitempty,oneof=todo in_progress completed blocked"`
	StartDate    *string    `json:"start_date,omitempty"`
	EndDate      *string    `json:"end_date,omitempty"`
}

// UpdateTaskRequest represents a request to update a task
type UpdateTaskRequest struct {
	Name         *string    `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description  *string    `json:"description,omitempty"`
	AssignedTo   *uuid.UUID `json:"assigned_to,omitempty"`
	PlannedHours *float64   `json:"planned_hours,omitempty" validate:"omitempty,min=0"`
	ActualHours  *float64   `json:"actual_hours,omitempty" validate:"omitempty,min=0"`
	Status       *string    `json:"status,omitempty" validate:"omitempty,oneof=todo in_progress completed blocked"`
	StartDate    *string    `json:"start_date,omitempty"`
	EndDate      *string    `json:"end_date,omitempty"`
}

// TaskResponse represents a task response
type TaskResponse struct {
	ID                 uuid.UUID             `json:"id"`
	ProjectID          uuid.UUID             `json:"project_id"`
	AssignedTo         *uuid.UUID            `json:"assigned_to,omitempty"`
	Name               string                `json:"name"`
	Description        *string               `json:"description,omitempty"`
	PlannedHours       float64               `json:"planned_hours"`
	ActualHours        float64               `json:"actual_hours"`
	VarianceHours      float64               `json:"variance_hours"`
	VariancePercentage float64               `json:"variance_percentage"`
	Status             string                `json:"status"`
	StartDate          *string               `json:"start_date,omitempty"`
	EndDate            *string               `json:"end_date,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	Assignee           *MemberBriefResponse  `json:"assignee,omitempty"`
}

// MemberBriefResponse represents a brief member response for nesting
type MemberBriefResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// TaskListResponse represents a paginated list of tasks
type TaskListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	Pagination Pagination     `json:"pagination"`
}

// ProjectSummaryResponse represents the project summary response
type ProjectSummaryResponse struct {
	ProjectID          uuid.UUID `json:"project_id"`
	TotalTasks         int       `json:"total_tasks"`
	TotalPlannedHours  float64   `json:"total_planned_hours"`
	TotalActualHours   float64   `json:"total_actual_hours"`
	VarianceHours      float64   `json:"variance_hours"`
	VariancePercentage float64   `json:"variance_percentage"`
	IsOverBudget       bool      `json:"is_over_budget"`
	CompletedTasks     int       `json:"completed_tasks"`
	InProgressTasks    int       `json:"in_progress_tasks"`
	TodoTasks          int       `json:"todo_tasks"`
	BlockedTasks       int       `json:"blocked_tasks"`
	CompletionRate     float64   `json:"completion_rate"`
}
