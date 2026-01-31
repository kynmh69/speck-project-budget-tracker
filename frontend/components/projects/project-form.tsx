'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { createProjectSchema, updateProjectSchema, CreateProjectFormData, UpdateProjectFormData } from '@/schemas/project-schema';
import { Project, projectStatusLabels, ProjectStatus } from '@/types/project';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';

interface ProjectFormProps {
  project?: Project;
  onSubmit: (data: CreateProjectFormData | UpdateProjectFormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export function ProjectForm({ project, onSubmit, onCancel, isLoading }: ProjectFormProps) {
  const isEditing = !!project;

  const form = useForm<CreateProjectFormData | UpdateProjectFormData>({
    resolver: zodResolver(isEditing ? updateProjectSchema : createProjectSchema),
    defaultValues: {
      name: project?.name || '',
      description: project?.description || '',
      status: project?.status || 'planning',
      budget_amount: project?.budget_amount || undefined,
      start_date: project?.start_date || '',
      end_date: project?.end_date || '',
    },
  });

  const handleSubmit = (data: CreateProjectFormData | UpdateProjectFormData) => {
    // Clean up empty strings
    const cleanedData = {
      ...data,
      description: data.description || undefined,
      start_date: data.start_date || undefined,
      end_date: data.end_date || undefined,
      budget_amount: data.budget_amount || undefined,
    };
    onSubmit(cleanedData);
  };

  const statusOptions: { value: ProjectStatus; label: string }[] = [
    { value: 'planning', label: projectStatusLabels.planning },
    { value: 'in_progress', label: projectStatusLabels.in_progress },
    { value: 'completed', label: projectStatusLabels.completed },
    { value: 'on_hold', label: projectStatusLabels.on_hold },
  ];

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>プロジェクト名 *</FormLabel>
              <FormControl>
                <Input placeholder="プロジェクト名を入力" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>説明</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="プロジェクトの説明を入力"
                  rows={4}
                  {...field}
                  value={field.value || ''}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>ステータス</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="ステータスを選択" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {statusOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="budget_amount"
          render={({ field }) => (
            <FormItem>
              <FormLabel>予算（円）</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  placeholder="予算を入力"
                  {...field}
                  value={field.value || ''}
                  onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : undefined)}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="grid grid-cols-2 gap-4">
          <FormField
            control={form.control}
            name="start_date"
            render={({ field }) => (
              <FormItem>
                <FormLabel>開始日</FormLabel>
                <FormControl>
                  <Input type="date" {...field} value={field.value || ''} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="end_date"
            render={({ field }) => (
              <FormItem>
                <FormLabel>終了日</FormLabel>
                <FormControl>
                  <Input type="date" {...field} value={field.value || ''} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <div className="flex justify-end space-x-4 pt-4">
          <Button type="button" variant="outline" onClick={onCancel}>
            キャンセル
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? '保存中...' : isEditing ? '更新' : '作成'}
          </Button>
        </div>
      </form>
    </Form>
  );
}
