import { useQuery } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import { DashboardResponse } from '@/types/dashboard';
import { ApiResponse } from '@/types/api';

// API functions
const dashboardApi = {
  // Get dashboard summary
  getDashboard: async (): Promise<DashboardResponse> => {
    const response = await apiClient.get<ApiResponse<DashboardResponse>>('/dashboard');
    return (response as any).data;
  },
};

// Query keys
export const dashboardKeys = {
  all: ['dashboard'] as const,
  summary: () => [...dashboardKeys.all, 'summary'] as const,
};

/**
 * Hook to fetch the dashboard summary (KPIs and recent projects)
 */
export function useDashboard() {
  return useQuery({
    queryKey: dashboardKeys.summary(),
    queryFn: dashboardApi.getDashboard,
  });
}
