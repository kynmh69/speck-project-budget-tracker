import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import {
  BudgetSummaryResponse,
  BudgetResponse,
  UpdateRevenueFormData,
  CreateTimeEntryFormData,
  UpdateTimeEntryFormData,
  TimeEntryResponse,
  TimeEntryListResponse,
} from '@/schemas/budget-schema';
import { ApiResponse } from '@/types/api';

// API functions
const budgetApi = {
  // Get budget summary for a project
  getBudgetSummary: async (projectId: string): Promise<BudgetSummaryResponse> => {
    const response = await apiClient.get<ApiResponse<BudgetSummaryResponse>>(
      `/projects/${projectId}/budget`
    );
    return (response as any).data;
  },

  // Update project revenue
  updateRevenue: async (projectId: string, data: UpdateRevenueFormData): Promise<BudgetResponse> => {
    const response = await apiClient.put<ApiResponse<BudgetResponse>>(
      `/projects/${projectId}/budget/revenue`,
      data
    );
    return (response as any).data;
  },

  // Create time entry
  createTimeEntry: async (data: CreateTimeEntryFormData): Promise<TimeEntryResponse> => {
    const response = await apiClient.post<ApiResponse<TimeEntryResponse>>('/time-entries', data);
    return (response as any).data;
  },

  // Get time entries
  listTimeEntries: async (params?: {
    project_id?: string;
    task_id?: string;
    member_id?: string;
    start_date?: string;
    end_date?: string;
    page?: number;
    per_page?: number;
  }): Promise<TimeEntryListResponse> => {
    const response = await apiClient.get<ApiResponse<TimeEntryListResponse>>('/time-entries', {
      params,
    });
    return (response as any).data;
  },

  // Get a single time entry
  getTimeEntry: async (entryId: string): Promise<TimeEntryResponse> => {
    const response = await apiClient.get<ApiResponse<TimeEntryResponse>>(
      `/time-entries/${entryId}`
    );
    return (response as any).data;
  },

  // Update time entry
  updateTimeEntry: async (
    entryId: string,
    data: UpdateTimeEntryFormData
  ): Promise<TimeEntryResponse> => {
    const response = await apiClient.put<ApiResponse<TimeEntryResponse>>(
      `/time-entries/${entryId}`,
      data
    );
    return (response as any).data;
  },

  // Delete time entry
  deleteTimeEntry: async (entryId: string): Promise<void> => {
    await apiClient.delete(`/time-entries/${entryId}`);
  },
};

// Query keys
export const budgetKeys = {
  all: ['budget'] as const,
  summaries: () => [...budgetKeys.all, 'summary'] as const,
  summary: (projectId: string) => [...budgetKeys.summaries(), projectId] as const,
  timeEntries: () => [...budgetKeys.all, 'time-entries'] as const,
  timeEntryList: (filters?: Record<string, unknown>) =>
    [...budgetKeys.timeEntries(), 'list', filters] as const,
  timeEntryDetail: (id: string) => [...budgetKeys.timeEntries(), 'detail', id] as const,
};

// Hooks

/**
 * Hook to fetch budget summary for a project
 */
export function useBudgetSummary(projectId: string) {
  return useQuery({
    queryKey: budgetKeys.summary(projectId),
    queryFn: () => budgetApi.getBudgetSummary(projectId),
    enabled: !!projectId,
  });
}

/**
 * Hook to update project revenue
 */
export function useUpdateRevenue() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ projectId, data }: { projectId: string; data: UpdateRevenueFormData }) =>
      budgetApi.updateRevenue(projectId, data),
    onSuccess: (_, { projectId }) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.summary(projectId) });
    },
  });
}

/**
 * Hook to fetch time entries
 */
export function useTimeEntries(params?: {
  project_id?: string;
  task_id?: string;
  member_id?: string;
  start_date?: string;
  end_date?: string;
  page?: number;
  per_page?: number;
}) {
  return useQuery({
    queryKey: budgetKeys.timeEntryList(params),
    queryFn: () => budgetApi.listTimeEntries(params),
  });
}

/**
 * Hook to fetch a single time entry
 */
export function useTimeEntry(entryId: string) {
  return useQuery({
    queryKey: budgetKeys.timeEntryDetail(entryId),
    queryFn: () => budgetApi.getTimeEntry(entryId),
    enabled: !!entryId,
  });
}

/**
 * Hook to create a time entry
 */
export function useCreateTimeEntry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTimeEntryFormData) => budgetApi.createTimeEntry(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.timeEntries() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.summaries() });
    },
  });
}

/**
 * Hook to update a time entry
 */
export function useUpdateTimeEntry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ entryId, data }: { entryId: string; data: UpdateTimeEntryFormData }) =>
      budgetApi.updateTimeEntry(entryId, data),
    onSuccess: (_, { entryId }) => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.timeEntryDetail(entryId) });
      queryClient.invalidateQueries({ queryKey: budgetKeys.timeEntries() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.summaries() });
    },
  });
}

/**
 * Hook to delete a time entry
 */
export function useDeleteTimeEntry() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (entryId: string) => budgetApi.deleteTimeEntry(entryId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: budgetKeys.timeEntries() });
      queryClient.invalidateQueries({ queryKey: budgetKeys.summaries() });
    },
  });
}
