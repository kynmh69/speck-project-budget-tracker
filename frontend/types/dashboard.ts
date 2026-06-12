import { Project } from './project';

// Dashboard summary response (GET /api/v1/dashboard)
export interface DashboardResponse {
  total_projects: number;
  active_projects: number;
  completed_projects: number;
  total_revenue: number;
  total_profit: number;
  average_profit_rate: number;
  recent_projects: Project[];
}
