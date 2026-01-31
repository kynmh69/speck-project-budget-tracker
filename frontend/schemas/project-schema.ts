import { z } from 'zod';

// Project status enum
export const projectStatusEnum = z.enum(['planning', 'in_progress', 'completed', 'on_hold']);

// Create project schema
export const createProjectSchema = z.object({
  name: z.string().min(1, 'プロジェクト名を入力してください').max(200, 'プロジェクト名は200文字以内で入力してください'),
  description: z.string().optional(),
  status: projectStatusEnum.optional().default('planning'),
  budget_amount: z.coerce.number().min(0, '予算は0以上を入力してください').optional().nullable(),
  start_date: z.string().optional().nullable(),
  end_date: z.string().optional().nullable(),
}).refine(
  (data) => {
    if (data.start_date && data.end_date) {
      return new Date(data.start_date) <= new Date(data.end_date);
    }
    return true;
  },
  {
    message: '終了日は開始日以降の日付を選択してください',
    path: ['end_date'],
  }
);

// Update project schema
export const updateProjectSchema = z.object({
  name: z.string().min(1, 'プロジェクト名を入力してください').max(200, 'プロジェクト名は200文字以内で入力してください').optional(),
  description: z.string().optional().nullable(),
  status: projectStatusEnum.optional(),
  budget_amount: z.coerce.number().min(0, '予算は0以上を入力してください').optional().nullable(),
  start_date: z.string().optional().nullable(),
  end_date: z.string().optional().nullable(),
});

// Types
export type ProjectStatus = z.infer<typeof projectStatusEnum>;
export type CreateProjectFormData = z.infer<typeof createProjectSchema>;
export type UpdateProjectFormData = z.infer<typeof updateProjectSchema>;
