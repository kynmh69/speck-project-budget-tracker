// Budget types

export interface Budget {
  id: string;
  project_id: string;
  revenue: number;
  total_cost: number;
  profit: number;
  profit_rate: number;
  currency: string;
  is_deficit: boolean;
}

// Cost breakdown
export interface CostBreakdown {
  labor_cost: number;
  total_hours: number;
  average_rate: number;
}

// Member cost
export interface MemberCost {
  member_id: string;
  member_name: string;
  hours: number;
  hourly_rate: number;
  cost: number;
  percentage: number;
}

// Task cost
export interface TaskCost {
  task_id: string;
  task_name: string;
  hours: number;
  cost: number;
}

// Budget summary response
export interface BudgetSummary {
  project_id: string;
  project_name: string;
  budget: Budget;
  cost_breakdown: CostBreakdown;
  member_costs: MemberCost[];
  task_costs?: TaskCost[];
  warning_message?: string;
}

// Update revenue request
export interface UpdateRevenueRequest {
  revenue: number;
  currency?: string;
}

// Time entry
export interface TimeEntry {
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

// Time entry summary
export interface TimeEntrySummary {
  total_hours: number;
  total_cost: number;
}

// Time entry list response
export interface TimeEntryListResponse {
  time_entries: TimeEntry[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
  summary?: TimeEntrySummary;
}

// Create time entry request
export interface CreateTimeEntryRequest {
  task_id: string;
  member_id: string;
  work_date: string;
  hours: number;
  comment?: string;
}

// Update time entry request
export interface UpdateTimeEntryRequest {
  work_date?: string;
  hours?: number;
  comment?: string;
}

// Currency display helper
export const formatCurrency = (amount: number, currency: string = 'JPY'): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: currency,
  }).format(amount);
};

// Profit status color
export const getProfitStatusColor = (profit: number): string => {
  if (profit > 0) return 'text-green-600';
  if (profit < 0) return 'text-red-600';
  return 'text-gray-600';
};

// Profit status badge color
export const getProfitBadgeColor = (isDeficit: boolean): string => {
  return isDeficit ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800';
};
