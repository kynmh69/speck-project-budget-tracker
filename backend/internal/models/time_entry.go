package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeEntry struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TaskID             uuid.UUID `gorm:"type:uuid;not null;index" json:"task_id"`
	MemberID           uuid.UUID `gorm:"type:uuid;not null;index" json:"member_id"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	WorkDate           time.Time `gorm:"type:date;not null;index" json:"work_date"`
	Hours              float64   `gorm:"type:decimal(5,2);not null" json:"hours"`
	HourlyRateSnapshot *float64  `gorm:"type:decimal(10,2)" json:"hourly_rate_snapshot,omitempty"`
	Comment            *string   `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Relations
	Task   Task   `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Member Member `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies table name
func (TimeEntry) TableName() string {
	return "time_entries"
}

// BeforeCreate hook
func (te *TimeEntry) BeforeCreate(tx *gorm.DB) error {
	if te.ID == uuid.Nil {
		te.ID = uuid.New()
	}
	return nil
}

// Cost calculates the cost of this time entry
func (te *TimeEntry) Cost() float64 {
	if te.HourlyRateSnapshot == nil {
		return 0
	}
	return te.Hours * *te.HourlyRateSnapshot
}
