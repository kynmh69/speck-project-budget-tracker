import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import {
  Project,
  ProjectDetail,
  ProjectListResponse,
  CreateProjectRequest,
  UpdateProjectRequest,
} from '@/types/project';
import { ApiResponse } from '@/types/api';

// API functions
const projectsApi = {
  // Create a new project
  create: async (data: CreateProjectRequest): Promise<Project> => {
    const response = await apiClient.post<ApiResponse<Project>>('/projects', data);
    return (response as any).data;
  },

  // Get projects list
  list: async (params?: {
    page?: number;
    per_page?: number;
    status?: string;
    search?: string;
    sort?: string;
    order?: string;
  }): Promise<ProjectListResponse> => {
    const response = await apiClient.get<ApiResponse<ProjectListResponse>>('/projects', { params });
    return (response as any).data;
  },

  // Get a single project
  get: async (projectId: string): Promise<ProjectDetail> => {
    const response = await apiClient.get<ApiResponse<ProjectDetail>>(`/projects/${projectId}`);
    return (response as any).data;
  },

  // Update a project
  update: async (projectId: string, data: UpdateProjectRequest): Promise<Project> => {
    const response = await apiClient.put<ApiResponse<Project>>(`/projects/${projectId}`, data);
    return (response as any).data;
  },

  // Delete a project
  delete: async (projectId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}`);
  },
};

// Query keys
export const projectKeys = {
  all: ['projects'] as const,
  lists: () => [...projectKeys.all, 'list'] as const,
  list: (filters?: Record<string, unknown>) => [...projectKeys.lists(), filters] as const,
  details: () => [...projectKeys.all, 'detail'] as const,
  detail: (id: string) => [...projectKeys.details(), id] as const,
};

// Hooks

/**
 * Hook to fetch projects list
 */
export function useProjects(params?: {
  page?: number;
  per_page?: number;
  status?: string;
  search?: string;
  sort?: string;
  order?: string;
}) {
  return useQuery({
    queryKey: projectKeys.list(params),
    queryFn: () => projectsApi.list(params),
  });
}

/**
 * Hook to fetch a single project
 */
export function useProject(projectId: string) {
  return useQuery({
    queryKey: projectKeys.detail(projectId),
    queryFn: () => projectsApi.get(projectId),
    enabled: !!projectId,
  });
}

/**
 * Hook to create a project
 */
export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProjectRequest) => projectsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
    },
  });
}

/**
 * Hook to update a project
 */
export function useUpdateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ projectId, data }: { projectId: string; data: UpdateProjectRequest }) =>
      projectsApi.update(projectId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
      queryClient.invalidateQueries({ queryKey: projectKeys.detail(variables.projectId) });
    },
  });
}

/**
 * Hook to delete a project
 */
export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (projectId: string) => projectsApi.delete(projectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
    },
  });
}
