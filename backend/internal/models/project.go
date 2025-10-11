package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string         `gorm:"type:varchar(200);not null;index" json:"name"`
	Description  *string        `gorm:"type:text" json:"description,omitempty"`
	Status       string         `gorm:"type:varchar(20);not null;default:'planning';index" json:"status"`
	BudgetAmount *float64       `gorm:"type:decimal(15,2)" json:"budget_amount,omitempty"`
	StartDate    *time.Time     `gorm:"type:date" json:"start_date,omitempty"`
	EndDate      *time.Time     `gorm:"type:date" json:"end_date,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tasks   []Task  `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	Budget  *Budget `gorm:"foreignKey:ProjectID" json:"budget,omitempty"`
	Members []Member `gorm:"many2many:project_members" json:"members,omitempty"`
}

// TableName specifies table name
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate hook
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
