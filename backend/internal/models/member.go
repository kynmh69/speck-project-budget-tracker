package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID     *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Email      string         `gorm:"type:varchar(255);not null;index" json:"email"`
	Role       *string        `gorm:"type:varchar(50)" json:"role,omitempty"`
	HourlyRate float64        `gorm:"type:decimal(10,2);default:0.00" json:"hourly_rate"`
	Department *string        `gorm:"type:varchar(100)" json:"department,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User        *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Projects    []Project     `gorm:"many2many:project_members" json:"projects,omitempty"`
	TimeEntries []TimeEntry   `gorm:"foreignKey:MemberID" json:"time_entries,omitempty"`
}

// TableName specifies table name
func (Member) TableName() string {
	return "members"
}

// BeforeCreate hook
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
