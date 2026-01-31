import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/lib/api-client';
import {
  MemberResponse,
  MemberListResponse,
  CreateMemberFormData,
  UpdateMemberFormData,
  AssignMemberFormData,
} from '@/schemas/member-schema';
import { ProjectMember } from '@/types/member';
import { ApiResponse } from '@/types/api';

// API functions
const membersApi = {
  // Create a new member
  create: async (data: CreateMemberFormData): Promise<MemberResponse> => {
    const response = await apiClient.post<ApiResponse<MemberResponse>>('/members', data);
    return (response as any).data;
  },

  // Get all members
  list: async (params?: {
    page?: number;
    per_page?: number;
    search?: string;
    department?: string;
  }): Promise<MemberListResponse> => {
    const response = await apiClient.get<ApiResponse<MemberListResponse>>('/members', { params });
    return (response as any).data;
  },

  // Get a single member
  get: async (memberId: string): Promise<MemberResponse> => {
    const response = await apiClient.get<ApiResponse<MemberResponse>>(`/members/${memberId}`);
    return (response as any).data;
  },

  // Update a member
  update: async (memberId: string, data: UpdateMemberFormData): Promise<MemberResponse> => {
    const response = await apiClient.put<ApiResponse<MemberResponse>>(`/members/${memberId}`, data);
    return (response as any).data;
  },

  // Delete a member
  delete: async (memberId: string): Promise<void> => {
    await apiClient.delete(`/members/${memberId}`);
  },

  // Get project members
  getProjectMembers: async (projectId: string): Promise<ProjectMember[]> => {
    const response = await apiClient.get<ApiResponse<ProjectMember[]>>(
      `/projects/${projectId}/members`
    );
    return (response as any).data;
  },

  // Assign member to project
  assignToProject: async (projectId: string, data: AssignMemberFormData): Promise<ProjectMember> => {
    const response = await apiClient.post<ApiResponse<ProjectMember>>(
      `/projects/${projectId}/members`,
      data
    );
    return (response as any).data;
  },

  // Remove member from project
  removeFromProject: async (projectId: string, memberId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/members/${memberId}`);
  },
};

// Query keys
export const memberKeys = {
  all: ['members'] as const,
  lists: () => [...memberKeys.all, 'list'] as const,
  list: (filters?: Record<string, unknown>) => [...memberKeys.lists(), filters] as const,
  details: () => [...memberKeys.all, 'detail'] as const,
  detail: (id: string) => [...memberKeys.details(), id] as const,
  projectMembers: (projectId: string) => [...memberKeys.all, 'project', projectId] as const,
};

// Hooks

/**
 * Hook to fetch all members
 */
export function useMembers(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  department?: string;
}) {
  return useQuery({
    queryKey: memberKeys.list(params),
    queryFn: () => membersApi.list(params),
  });
}

/**
 * Hook to fetch a single member
 */
export function useMember(memberId: string) {
  return useQuery({
    queryKey: memberKeys.detail(memberId),
    queryFn: () => membersApi.get(memberId),
    enabled: !!memberId,
  });
}

/**
 * Hook to fetch project members
 */
export function useProjectMembers(projectId: string) {
  return useQuery({
    queryKey: memberKeys.projectMembers(projectId),
    queryFn: () => membersApi.getProjectMembers(projectId),
    enabled: !!projectId,
  });
}

/**
 * Hook to create a member
 */
export function useCreateMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateMemberFormData) => membersApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: memberKeys.lists() });
    },
  });
}

/**
 * Hook to update a member
 */
export function useUpdateMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ memberId, data }: { memberId: string; data: UpdateMemberFormData }) =>
      membersApi.update(memberId, data),
    onSuccess: (_, { memberId }) => {
      queryClient.invalidateQueries({ queryKey: memberKeys.detail(memberId) });
      queryClient.invalidateQueries({ queryKey: memberKeys.lists() });
    },
  });
}

/**
 * Hook to delete a member
 */
export function useDeleteMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (memberId: string) => membersApi.delete(memberId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: memberKeys.lists() });
    },
  });
}

/**
 * Hook to assign a member to a project
 */
export function useAssignMemberToProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ projectId, data }: { projectId: string; data: AssignMemberFormData }) =>
      membersApi.assignToProject(projectId, data),
    onSuccess: (_, { projectId }) => {
      queryClient.invalidateQueries({ queryKey: memberKeys.projectMembers(projectId) });
    },
  });
}

/**
 * Hook to remove a member from a project
 */
export function useRemoveMemberFromProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ projectId, memberId }: { projectId: string; memberId: string }) =>
      membersApi.removeFromProject(projectId, memberId),
    onSuccess: (_, { projectId }) => {
      queryClient.invalidateQueries({ queryKey: memberKeys.projectMembers(projectId) });
    },
  });
}
