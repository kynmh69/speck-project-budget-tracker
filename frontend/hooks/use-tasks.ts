import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import {
  TaskResponse,
  TaskListResponse,
  ProjectSummary,
  CreateTaskFormData,
  UpdateTaskFormData,
} from '@/schemas/task-schema';
import { ApiResponse } from '@/types/api';

// API functions
const tasksApi = {
  // Create a new task
  create: async (projectId: string, data: CreateTaskFormData): Promise<TaskResponse> => {
    const response = await apiClient.post<ApiResponse<TaskResponse>>(
      `/projects/${projectId}/tasks`,
      data
    );
    return (response as any).data;
  },

  // Get tasks for a project
  list: async (
    projectId: string,
    params?: { page?: number; per_page?: number; status?: string }
  ): Promise<TaskListResponse> => {
    const response = await apiClient.get<ApiResponse<TaskListResponse>>(
      `/projects/${projectId}/tasks`,
      { params }
    );
    return (response as any).data;
  },

  // Get a single task
  get: async (taskId: string): Promise<TaskResponse> => {
    const response = await apiClient.get<ApiResponse<TaskResponse>>(`/tasks/${taskId}`);
    return (response as any).data;
  },

  // Update a task
  update: async (taskId: string, data: UpdateTaskFormData): Promise<TaskResponse> => {
    const response = await apiClient.put<ApiResponse<TaskResponse>>(`/tasks/${taskId}`, data);
    return (response as any).data;
  },

  // Delete a task
  delete: async (taskId: string): Promise<void> => {
    await apiClient.delete(`/tasks/${taskId}`);
  },

  // Get project summary
  getSummary: async (projectId: string): Promise<ProjectSummary> => {
    const response = await apiClient.get<ApiResponse<ProjectSummary>>(
      `/projects/${projectId}/summary`
    );
    return (response as any).data;
  },
};

// Query keys
export const taskKeys = {
  all: ['tasks'] as const,
  lists: () => [...taskKeys.all, 'list'] as const,
  list: (projectId: string, filters?: Record<string, unknown>) =>
    [...taskKeys.lists(), projectId, filters] as const,
  details: () => [...taskKeys.all, 'detail'] as const,
  detail: (id: string) => [...taskKeys.details(), id] as const,
  summaries: () => [...taskKeys.all, 'summary'] as const,
  summary: (projectId: string) => [...taskKeys.summaries(), projectId] as const,
};

// Hooks

/**
 * Hook to fetch tasks for a project
 */
export function useTasks(
  projectId: string,
  params?: { page?: number; per_page?: number; status?: string }
) {
  return useQuery({
    queryKey: taskKeys.list(projectId, params),
    queryFn: () => tasksApi.list(projectId, params),
    enabled: !!projectId,
  });
}

/**
 * Hook to fetch a single task
 */
export function useTask(taskId: string) {
  return useQuery({
    queryKey: taskKeys.detail(taskId),
    queryFn: () => tasksApi.get(taskId),
    enabled: !!taskId,
  });
}

/**
 * Hook to fetch project summary
 */
export function useProjectSummary(projectId: string) {
  return useQuery({
    queryKey: taskKeys.summary(projectId),
    queryFn: () => tasksApi.getSummary(projectId),
    enabled: !!projectId,
  });
}

/**
 * Hook to create a task
 */
export function useCreateTask(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTaskFormData) => tasksApi.create(projectId, data),
    onSuccess: () => {
      // Invalidate task list and project summary
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.summary(projectId) });
    },
  });
}

/**
 * Hook to update a task
 */
export function useUpdateTask(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ taskId, data }: { taskId: string; data: UpdateTaskFormData }) =>
      tasksApi.update(taskId, data),
    onSuccess: (_, variables) => {
      // Invalidate specific task and list
      queryClient.invalidateQueries({ queryKey: taskKeys.detail(variables.taskId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.summary(projectId) });
    },
  });
}

/**
 * Hook to delete a task
 */
export function useDeleteTask(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (taskId: string) => tasksApi.delete(taskId),
    onSuccess: () => {
      // Invalidate task list and project summary
      queryClient.invalidateQueries({ queryKey: taskKeys.list(projectId) });
      queryClient.invalidateQueries({ queryKey: taskKeys.summary(projectId) });
    },
  });
}

/**
 * Hook to update task hours (convenience wrapper)
 */
export function useUpdateTaskHours(projectId: string) {
  const updateTask = useUpdateTask(projectId);

  return {
    ...updateTask,
    updatePlannedHours: (taskId: string, plannedHours: number) =>
      updateTask.mutate({ taskId, data: { planned_hours: plannedHours } }),
    updateActualHours: (taskId: string, actualHours: number) =>
      updateTask.mutate({ taskId, data: { actual_hours: actualHours } }),
  };
}
