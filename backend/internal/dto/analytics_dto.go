package dto

// PlanActualTaskData represents plan vs actual hours for a single task
type PlanActualTaskData struct {
	TaskName     string  `json:"task_name"`
	PlannedHours float64 `json:"planned_hours"`
	ActualHours  float64 `json:"actual_hours"`
	Variance     float64 `json:"variance"`
}

// PlanActualResponse represents plan vs actual comparison data
type PlanActualResponse struct {
	Tasks []PlanActualTaskData `json:"tasks"`
}

// MemberCostData represents cost attributed to a member
type MemberCostData struct {
	MemberName string  `json:"member_name"`
	Cost       float64 `json:"cost"`
}

// BudgetAnalyticsResponse represents budget analytics data
type BudgetAnalyticsResponse struct {
	Revenue       float64          `json:"revenue"`
	TotalCost     float64          `json:"total_cost"`
	CostBreakdown []MemberCostData `json:"cost_breakdown"`
	Profit        float64          `json:"profit"`
	ProfitRate    float64          `json:"profit_rate"`
}

// TrendDataPoint represents one bucket in a trend series
type TrendDataPoint struct {
	Date         string  `json:"date"`
	PlannedHours float64 `json:"planned_hours"`
	ActualHours  float64 `json:"actual_hours"`
	Cost         float64 `json:"cost"`
}

// TrendResponse represents time-series trend data
type TrendResponse struct {
	Period     string           `json:"period"`
	DataPoints []TrendDataPoint `json:"data_points"`
}

// TaskDistributionItem represents a task's share of actual hours
type TaskDistributionItem struct {
	TaskName    string  `json:"task_name"`
	ActualHours float64 `json:"actual_hours"`
	Percentage  float64 `json:"percentage"`
}

// TaskDistributionResponse represents actual hours distribution across tasks
type TaskDistributionResponse struct {
	Tasks []TaskDistributionItem `json:"tasks"`
}

// ProjectComparisonItem represents one project's KPIs for comparison
type ProjectComparisonItem struct {
	ProjectID    string  `json:"project_id"`
	ProjectName  string  `json:"project_name"`
	PlannedHours float64 `json:"planned_hours"`
	ActualHours  float64 `json:"actual_hours"`
	Revenue      float64 `json:"revenue"`
	TotalCost    float64 `json:"total_cost"`
	Profit       float64 `json:"profit"`
	ProfitRate   float64 `json:"profit_rate"`
}

// ProjectsComparisonResponse represents cross-project comparison data
type ProjectsComparisonResponse struct {
	Projects []ProjectComparisonItem `json:"projects"`
}
