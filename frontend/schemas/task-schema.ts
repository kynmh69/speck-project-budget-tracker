import { z } from 'zod';

// Task status enum
export const taskStatusEnum = z.enum(['todo', 'in_progress', 'completed', 'blocked']);

// Create task schema
export const createTaskSchema = z.object({
  name: z.string().min(1, 'タスク名を入力してください').max(200, 'タスク名は200文字以内で入力してください'),
  description: z.string().optional(),
  assigned_to: z.string().uuid().optional().nullable(),
  planned_hours: z.coerce.number().min(0, '予定工数は0以上を入力してください').default(0),
  status: taskStatusEnum.optional().default('todo'),
  start_date: z.string().optional().nullable(),
  end_date: z.string().optional().nullable(),
});

// Update task schema
export const updateTaskSchema = z.object({
  name: z.string().min(1, 'タスク名を入力してください').max(200, 'タスク名は200文字以内で入力してください').optional(),
  description: z.string().optional().nullable(),
  assigned_to: z.string().uuid().optional().nullable(),
  planned_hours: z.coerce.number().min(0, '予定工数は0以上を入力してください').optional(),
  actual_hours: z.coerce.number().min(0, '実績工数は0以上を入力してください').optional(),
  status: taskStatusEnum.optional(),
  start_date: z.string().optional().nullable(),
  end_date: z.string().optional().nullable(),
});

// Types
export type TaskStatus = z.infer<typeof taskStatusEnum>;
export type CreateTaskFormData = z.infer<typeof createTaskSchema>;
export type UpdateTaskFormData = z.infer<typeof updateTaskSchema>;

// Task response type (from API)
export interface TaskResponse {
  id: string;
  project_id: string;
  assigned_to?: string;
  name: string;
  description?: string;
  planned_hours: number;
  actual_hours: number;
  variance_hours: number;
  variance_percentage: number;
  status: TaskStatus;
  start_date?: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
  assignee?: {
    id: string;
    name: string;
  };
}

// Project summary type
export interface ProjectSummary {
  project_id: string;
  total_tasks: number;
  total_planned_hours: number;
  total_actual_hours: number;
  variance_hours: number;
  variance_percentage: number;
  is_over_budget: boolean;
  completed_tasks: number;
  in_progress_tasks: number;
  todo_tasks: number;
  blocked_tasks: number;
  completion_rate: number;
}

// Task list response
export interface TaskListResponse {
  tasks: TaskResponse[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

// Status label mapping
export const taskStatusLabels: Record<TaskStatus, string> = {
  todo: '未着手',
  in_progress: '進行中',
  completed: '完了',
  blocked: 'ブロック中',
};

// Status color mapping
export const taskStatusColors: Record<TaskStatus, string> = {
  todo: 'bg-gray-100 text-gray-800',
  in_progress: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  blocked: 'bg-red-100 text-red-800',
};
