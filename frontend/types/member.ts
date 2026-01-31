// Member types

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

// Member list response
export interface MemberListResponse {
  members: Member[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

// Create member request
export interface CreateMemberRequest {
  name: string;
  email: string;
  role?: string;
  hourly_rate: number;
  department?: string;
  user_id?: string;
}

// Update member request
export interface UpdateMemberRequest {
  name?: string;
  email?: string;
  role?: string;
  hourly_rate?: number;
  department?: string;
}

// Project member assignment
export interface ProjectMember {
  id: string;
  project_id: string;
  member_id: string;
  joined_at: string;
  left_at?: string;
  allocation_rate: number;
  hourly_rate_snapshot?: number;
  member?: Member;
}

// Assign member request
export interface AssignMemberRequest {
  member_id: string;
  allocation_rate?: number;
  hourly_rate_snapshot?: number;
}
