import { z } from 'zod';

// Update revenue schema
export const updateRevenueSchema = z.object({
  revenue: z.number().min(0, '売上は0以上で入力してください'),
  currency: z.string().length(3).optional(),
});

// Create time entry schema
export const createTimeEntrySchema = z.object({
  task_id: z.string().uuid('タスクを選択してください'),
  member_id: z.string().uuid('メンバーを選択してください'),
  work_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, '有効な日付を入力してください'),
  hours: z.number().min(0, '工数は0以上で入力してください').max(24, '工数は24時間以内で入力してください'),
  comment: z.string().optional(),
});

// Update time entry schema
export const updateTimeEntrySchema = z.object({
  work_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/).optional(),
  hours: z.number().min(0).max(24).optional(),
  comment: z.string().optional().nullable(),
});

// Types inferred from schemas
export type UpdateRevenueFormData = z.infer<typeof updateRevenueSchema>;
export type CreateTimeEntryFormData = z.infer<typeof createTimeEntrySchema>;
export type UpdateTimeEntryFormData = z.infer<typeof updateTimeEntrySchema>;

// Budget response types
export interface BudgetResponse {
  id: string;
  project_id: string;
  revenue: number;
  total_cost: number;
  profit: number;
  profit_rate: number;
  currency: string;
  is_deficit: boolean;
}

export interface CostBreakdownResponse {
  labor_cost: number;
  total_hours: number;
  average_rate: number;
}

export interface MemberCostResponse {
  member_id: string;
  member_name: string;
  hours: number;
  hourly_rate: number;
  cost: number;
  percentage: number;
}

export interface BudgetSummaryResponse {
  project_id: string;
  project_name: string;
  budget: BudgetResponse;
  cost_breakdown: CostBreakdownResponse;
  member_costs: MemberCostResponse[];
  warning_message?: string;
}

export interface TimeEntryResponse {
  id: string;
  task_id: string;
  member_id: string;
  user_id: string;
  work_date: string;
  hours: number;
  hourly_rate_snapshot?: number;
  cost: number;
  comment?: string;
  created_at: string;
  updated_at: string;
  member?: {
    id: string;
    name: string;
  };
}

export interface TimeEntrySummary {
  total_hours: number;
  total_cost: number;
}

export interface TimeEntryListResponse {
  time_entries: TimeEntryResponse[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
  summary?: TimeEntrySummary;
}
