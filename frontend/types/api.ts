export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

export interface Project {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  status: 'planning' | 'in_progress' | 'completed' | 'on_hold';
  budget_amount?: number;
  start_date?: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  project_id: string;
  assigned_to?: string;
  name: string;
  description?: string;
  planned_hours: number;
  actual_hours: number;
  status: 'todo' | 'in_progress' | 'completed' | 'blocked';
  start_date?: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Member {
  id: string;
  user_id?: string;
  name: string;
  email: string;
  role?: string;
  hourly_rate: number;
  department?: string;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: string;
  project_id: string;
  revenue: number;
  total_cost: number;
  profit: number;
  profit_rate: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface TimeEntry {
  id: string;
  task_id: string;
  member_id: string;
  user_id: string;
  work_date: string;
  hours: number;
  hourly_rate_snapshot?: number;
  comment?: string;
  created_at: string;
}

export interface Pagination {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
}

export interface ListResponse<T> {
  data: T[];
  pagination: Pagination;
}
