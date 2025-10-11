package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectMember struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProjectID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"project_id"`
	MemberID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"member_id"`
	JoinedAt           time.Time  `gorm:"type:date;not null;default:CURRENT_DATE" json:"joined_at"`
	LeftAt             *time.Time `gorm:"type:date" json:"left_at,omitempty"`
	AllocationRate     float64    `gorm:"type:decimal(3,2);default:1.00" json:"allocation_rate"`
	HourlyRateSnapshot *float64   `gorm:"type:decimal(10,2)" json:"hourly_rate_snapshot,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Relations
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Member  Member  `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

// TableName specifies table name
func (ProjectMember) TableName() string {
	return "project_members"
}

// BeforeCreate hook
func (pm *ProjectMember) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}

// IsActive checks if the member is currently active on the project
func (pm *ProjectMember) IsActive() bool {
	return pm.LeftAt == nil || pm.LeftAt.After(time.Now())
}
