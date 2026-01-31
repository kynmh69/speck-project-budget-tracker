import { z } from 'zod';

// Create member schema
export const createMemberSchema = z.object({
  name: z.string().min(1, '名前は必須です').max(100, '名前は100文字以内で入力してください'),
  email: z.string().email('有効なメールアドレスを入力してください').max(255),
  role: z.string().max(50).optional(),
  hourly_rate: z.number().min(0, '時給は0以上で入力してください'),
  department: z.string().max(100).optional(),
  user_id: z.string().uuid().optional(),
});

// Update member schema
export const updateMemberSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  email: z.string().email().max(255).optional(),
  role: z.string().max(50).optional().nullable(),
  hourly_rate: z.number().min(0).optional(),
  department: z.string().max(100).optional().nullable(),
});

// Assign member to project schema
export const assignMemberSchema = z.object({
  member_id: z.string().uuid('有効なメンバーIDを選択してください'),
  allocation_rate: z.number().min(0).max(1).optional(),
  hourly_rate_snapshot: z.number().min(0).optional(),
});

// Types inferred from schemas
export type CreateMemberFormData = z.infer<typeof createMemberSchema>;
export type UpdateMemberFormData = z.infer<typeof updateMemberSchema>;
export type AssignMemberFormData = z.infer<typeof assignMemberSchema>;

// Member response type
export interface MemberResponse {
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
  members: MemberResponse[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}
