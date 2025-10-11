package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"project_id"`
	Revenue    float64   `gorm:"type:decimal(15,2);default:0.00" json:"revenue"`
	TotalCost  float64   `gorm:"type:decimal(15,2);default:0.00" json:"total_cost"`
	Profit     float64   `gorm:"type:decimal(15,2);default:0.00" json:"profit"`
	ProfitRate float64   `gorm:"type:decimal(5,2);default:0.00" json:"profit_rate"`
	Currency   string    `gorm:"type:varchar(3);not null;default:'JPY'" json:"currency"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName specifies table name
func (Budget) TableName() string {
	return "budgets"
}

// BeforeCreate hook
func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// CalculateProfit calculates profit and profit rate
func (b *Budget) CalculateProfit() {
	b.Profit = b.Revenue - b.TotalCost
	if b.Revenue > 0 {
		b.ProfitRate = (b.Profit / b.Revenue) * 100
	} else {
		b.ProfitRate = 0
	}
}
