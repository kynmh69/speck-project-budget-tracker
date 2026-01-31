package dto

import (
	"github.com/google/uuid"
)

// UpdateRevenueRequest represents a request to update project revenue
type UpdateRevenueRequest struct {
	Revenue  float64 `json:"revenue" validate:"min=0"`
	Currency *string `json:"currency,omitempty" validate:"omitempty,len=3"`
}

// BudgetResponse represents a budget response
type BudgetResponse struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Revenue    float64   `json:"revenue"`
	TotalCost  float64   `json:"total_cost"`
	Profit     float64   `json:"profit"`
	ProfitRate float64   `json:"profit_rate"`
	Currency   string    `json:"currency"`
	IsDeficit  bool      `json:"is_deficit"`
}

// BudgetSummaryResponse represents a comprehensive budget summary
type BudgetSummaryResponse struct {
	ProjectID      uuid.UUID               `json:"project_id"`
	ProjectName    string                  `json:"project_name"`
	Budget         BudgetResponse          `json:"budget"`
	CostBreakdown  CostBreakdownResponse   `json:"cost_breakdown"`
	MemberCosts    []MemberCostResponse    `json:"member_costs"`
	TaskCosts      []TaskCostResponse      `json:"task_costs,omitempty"`
	WarningMessage *string                 `json:"warning_message,omitempty"`
}

// CostBreakdownResponse represents cost breakdown by category
type CostBreakdownResponse struct {
	LaborCost     float64 `json:"labor_cost"`
	TotalHours    float64 `json:"total_hours"`
	AverageRate   float64 `json:"average_rate"`
}

// MemberCostResponse represents cost breakdown by member
type MemberCostResponse struct {
	MemberID   uuid.UUID `json:"member_id"`
	MemberName string    `json:"member_name"`
	Hours      float64   `json:"hours"`
	HourlyRate float64   `json:"hourly_rate"`
	Cost       float64   `json:"cost"`
	Percentage float64   `json:"percentage"`
}

// TaskCostResponse represents cost breakdown by task
type TaskCostResponse struct {
	TaskID   uuid.UUID `json:"task_id"`
	TaskName string    `json:"task_name"`
	Hours    float64   `json:"hours"`
	Cost     float64   `json:"cost"`
}

// BudgetComparisonResponse represents budget comparison between planned and actual
type BudgetComparisonResponse struct {
	ProjectID     uuid.UUID `json:"project_id"`
	PlannedBudget float64   `json:"planned_budget"`
	ActualCost    float64   `json:"actual_cost"`
	Variance      float64   `json:"variance"`
	VarianceRate  float64   `json:"variance_rate"`
	IsOverBudget  bool      `json:"is_over_budget"`
}
