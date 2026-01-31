package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateMemberRequest represents a request to create a member
type CreateMemberRequest struct {
	Name       string     `json:"name" validate:"required,min=1,max=100"`
	Email      string     `json:"email" validate:"required,email,max=255"`
	Role       *string    `json:"role,omitempty" validate:"omitempty,max=50"`
	HourlyRate float64    `json:"hourly_rate" validate:"min=0"`
	Department *string    `json:"department,omitempty" validate:"omitempty,max=100"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
}

// UpdateMemberRequest represents a request to update a member
type UpdateMemberRequest struct {
	Name       *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Email      *string  `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Role       *string  `json:"role,omitempty" validate:"omitempty,max=50"`
	HourlyRate *float64 `json:"hourly_rate,omitempty" validate:"omitempty,min=0"`
	Department *string  `json:"department,omitempty" validate:"omitempty,max=100"`
}

// MemberResponse represents a member response
type MemberResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Role       *string    `json:"role,omitempty"`
	HourlyRate float64    `json:"hourly_rate"`
	Department *string    `json:"department,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// MemberListResponse represents a paginated list of members
type MemberListResponse struct {
	Members    []MemberResponse `json:"members"`
	Pagination Pagination       `json:"pagination"`
}

// AssignMemberRequest represents a request to assign a member to a project
type AssignMemberRequest struct {
	MemberID           uuid.UUID `json:"member_id" validate:"required"`
	AllocationRate     *float64  `json:"allocation_rate,omitempty" validate:"omitempty,min=0,max=1"`
	HourlyRateSnapshot *float64  `json:"hourly_rate_snapshot,omitempty" validate:"omitempty,min=0"`
}

// ProjectMemberResponse represents a project member assignment response
type ProjectMemberResponse struct {
	ID                 uuid.UUID       `json:"id"`
	ProjectID          uuid.UUID       `json:"project_id"`
	MemberID           uuid.UUID       `json:"member_id"`
	JoinedAt           string          `json:"joined_at"`
	LeftAt             *string         `json:"left_at,omitempty"`
	AllocationRate     float64         `json:"allocation_rate"`
	HourlyRateSnapshot *float64        `json:"hourly_rate_snapshot,omitempty"`
	Member             *MemberResponse `json:"member,omitempty"`
}

// CreateTimeEntryRequest represents a request to create a time entry
type CreateTimeEntryRequest struct {
	TaskID   uuid.UUID `json:"task_id" validate:"required"`
	MemberID uuid.UUID `json:"member_id" validate:"required"`
	WorkDate string    `json:"work_date" validate:"required"`
	Hours    float64   `json:"hours" validate:"required,min=0,max=24"`
	Comment  *string   `json:"comment,omitempty"`
}

// UpdateTimeEntryRequest represents a request to update a time entry
type UpdateTimeEntryRequest struct {
	WorkDate *string  `json:"work_date,omitempty"`
	Hours    *float64 `json:"hours,omitempty" validate:"omitempty,min=0,max=24"`
	Comment  *string  `json:"comment,omitempty"`
}

// TimeEntryResponse represents a time entry response
type TimeEntryResponse struct {
	ID                 uuid.UUID            `json:"id"`
	TaskID             uuid.UUID            `json:"task_id"`
	MemberID           uuid.UUID            `json:"member_id"`
	UserID             uuid.UUID            `json:"user_id"`
	WorkDate           string               `json:"work_date"`
	Hours              float64              `json:"hours"`
	HourlyRateSnapshot *float64             `json:"hourly_rate_snapshot,omitempty"`
	Cost               float64              `json:"cost"`
	Comment            *string              `json:"comment,omitempty"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	Member             *MemberBriefResponse `json:"member,omitempty"`
}

// TimeEntryListResponse represents a paginated list of time entries
type TimeEntryListResponse struct {
	TimeEntries []TimeEntryResponse `json:"time_entries"`
	Pagination  Pagination          `json:"pagination"`
	Summary     *TimeEntrySummary   `json:"summary,omitempty"`
}

// TimeEntrySummary represents the summary of time entries
type TimeEntrySummary struct {
	TotalHours float64 `json:"total_hours"`
	TotalCost  float64 `json:"total_cost"`
}
