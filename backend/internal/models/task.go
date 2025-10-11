package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProjectID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	AssignedTo   *uuid.UUID     `gorm:"type:uuid;index" json:"assigned_to,omitempty"`
	Name         string         `gorm:"type:varchar(200);not null" json:"name"`
	Description  *string        `gorm:"type:text" json:"description,omitempty"`
	PlannedHours float64        `gorm:"type:decimal(10,2);default:0.00" json:"planned_hours"`
	ActualHours  float64        `gorm:"type:decimal(10,2);default:0.00" json:"actual_hours"`
	Status       string         `gorm:"type:varchar(20);not null;default:'todo';index" json:"status"`
	StartDate    *time.Time     `gorm:"type:date" json:"start_date,omitempty"`
	EndDate      *time.Time     `gorm:"type:date" json:"end_date,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Project     Project     `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Assignee    *Member     `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	TimeEntries []TimeEntry `gorm:"foreignKey:TaskID" json:"time_entries,omitempty"`
}

// TableName specifies table name
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate hook
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// VarianceHours calculates the variance between planned and actual hours
func (t *Task) VarianceHours() float64 {
	return t.ActualHours - t.PlannedHours
}

// VariancePercentage calculates the variance percentage
func (t *Task) VariancePercentage() float64 {
	if t.PlannedHours == 0 {
		return 0
	}
	return (t.VarianceHours() / t.PlannedHours) * 100
}
