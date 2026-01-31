// Project types
export interface Project {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  status: ProjectStatus;
  budget_amount?: number;
  start_date?: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export type ProjectStatus = 'planning' | 'in_progress' | 'completed' | 'on_hold';

// Project detail with stats
export interface ProjectDetail extends Project {
  stats?: ProjectStats;
}

// Project statistics
export interface ProjectStats {
  total_tasks: number;
  completed_tasks: number;
  total_planned_hours: number;
  total_actual_hours: number;
  completion_rate: number;
}

// Project list response
export interface ProjectListResponse {
  projects: Project[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

// Create project request
export interface CreateProjectRequest {
  name: string;
  description?: string;
  status?: ProjectStatus;
  budget_amount?: number | null;
  start_date?: string | null;
  end_date?: string | null;
}

// Update project request
export interface UpdateProjectRequest {
  name?: string;
  description?: string | null;
  status?: ProjectStatus;
  budget_amount?: number | null;
  start_date?: string | null;
  end_date?: string | null;
}

// Status label mapping
export const projectStatusLabels: Record<ProjectStatus, string> = {
  planning: '計画中',
  in_progress: '進行中',
  completed: '完了',
  on_hold: '保留中',
};

// Status color mapping
export const projectStatusColors: Record<ProjectStatus, string> = {
  planning: 'bg-yellow-100 text-yellow-800',
  in_progress: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  on_hold: 'bg-gray-100 text-gray-800',
};
