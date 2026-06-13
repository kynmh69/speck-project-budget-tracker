// Analytics API (GET /api/v1/projects/{id}/analytics/*, /api/v1/analytics/*) types

export interface PlanActualTaskData {
  task_name: string;
  planned_hours: number;
  actual_hours: number;
  variance: number;
}

export interface PlanActualResponse {
  tasks: PlanActualTaskData[];
}

export interface MemberCostData {
  member_name: string;
  cost: number;
}

export interface BudgetAnalyticsResponse {
  revenue: number;
  total_cost: number;
  cost_breakdown: MemberCostData[];
  profit: number;
  profit_rate: number;
}

export type TrendPeriod = 'daily' | 'weekly' | 'monthly';

export interface TrendDataPoint {
  date: string;
  planned_hours: number;
  actual_hours: number;
  cost: number;
}

export interface TrendResponse {
  period: TrendPeriod;
  data_points: TrendDataPoint[];
}

export interface TaskDistributionItem {
  task_name: string;
  actual_hours: number;
  percentage: number;
}

export interface TaskDistributionResponse {
  tasks: TaskDistributionItem[];
}

export interface ProjectComparisonItem {
  project_id: string;
  project_name: string;
  planned_hours: number;
  actual_hours: number;
  revenue: number;
  total_cost: number;
  profit: number;
  profit_rate: number;
}

export interface ProjectsComparisonResponse {
  projects: ProjectComparisonItem[];
}
