package dto

// DashboardResponse represents the dashboard summary response
type DashboardResponse struct {
	TotalProjects     int64             `json:"total_projects"`
	ActiveProjects    int64             `json:"active_projects"`
	CompletedProjects int64             `json:"completed_projects"`
	TotalRevenue      float64           `json:"total_revenue"`
	TotalProfit       float64           `json:"total_profit"`
	AverageProfitRate float64           `json:"average_profit_rate"`
	RecentProjects    []ProjectResponse `json:"recent_projects"`
}
