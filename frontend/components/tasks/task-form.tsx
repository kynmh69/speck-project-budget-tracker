'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  createTaskSchema,
  CreateTaskFormData,
  UpdateTaskFormData,
  TaskResponse,
  taskStatusLabels,
} from '@/schemas/task-schema';

interface TaskFormProps {
  task?: TaskResponse;
  onSubmit: (data: CreateTaskFormData | UpdateTaskFormData) => void;
  onCancel: () => void;
  isLoading?: boolean;
}

export function TaskForm({ task, onSubmit, onCancel, isLoading }: TaskFormProps) {
  const isEditing = !!task;
  
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateTaskFormData>({
    resolver: zodResolver(createTaskSchema),
    defaultValues: {
      name: task?.name || '',
      description: task?.description || '',
      planned_hours: task?.planned_hours || 0,
      status: task?.status || 'todo',
      start_date: task?.start_date || '',
      end_date: task?.end_date || '',
    },
  });

  const handleFormSubmit = (data: CreateTaskFormData) => {
    if (isEditing) {
      // For update, include actual_hours from the form
      const formElement = document.getElementById('actual_hours') as HTMLInputElement;
      const actualHours = formElement ? parseFloat(formElement.value) || 0 : task?.actual_hours || 0;
      onSubmit({ ...data, actual_hours: actualHours } as UpdateTaskFormData);
    } else {
      onSubmit(data);
    }
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {/* タスク名 */}
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700">
          タスク名 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="name"
          {...register('name')}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          placeholder="タスク名を入力"
        />
        {errors.name && (
          <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
        )}
      </div>

      {/* 説明 */}
      <div>
        <label htmlFor="description" className="block text-sm font-medium text-gray-700">
          説明
        </label>
        <textarea
          id="description"
          {...register('description')}
          rows={3}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          placeholder="タスクの説明を入力"
        />
      </div>

      {/* ステータス */}
      <div>
        <label htmlFor="status" className="block text-sm font-medium text-gray-700">
          ステータス
        </label>
        <select
          id="status"
          {...register('status')}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        >
          {Object.entries(taskStatusLabels).map(([value, label]) => (
            <option key={value} value={value}>
              {label}
            </option>
          ))}
        </select>
      </div>

      {/* 工数入力 */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label htmlFor="planned_hours" className="block text-sm font-medium text-gray-700">
            予定工数（時間）
          </label>
          <input
            type="number"
            id="planned_hours"
            step="0.5"
            min="0"
            {...register('planned_hours')}
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
          {errors.planned_hours && (
            <p className="mt-1 text-sm text-red-600">{errors.planned_hours.message}</p>
          )}
        </div>

        {isEditing && (
          <div>
            <label htmlFor="actual_hours" className="block text-sm font-medium text-gray-700">
              実績工数（時間）
            </label>
            <input
              type="number"
              id="actual_hours"
              name="actual_hours"
              step="0.5"
              min="0"
              defaultValue={task?.actual_hours || 0}
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
        )}
      </div>

      {/* 日付入力 */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label htmlFor="start_date" className="block text-sm font-medium text-gray-700">
            開始日
          </label>
          <input
            type="date"
            id="start_date"
            {...register('start_date')}
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
        </div>

        <div>
          <label htmlFor="end_date" className="block text-sm font-medium text-gray-700">
            終了日
          </label>
          <input
            type="date"
            id="end_date"
            {...register('end_date')}
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* ボタン */}
      <div className="flex justify-end gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          キャンセル
        </button>
        <button
          type="submit"
          disabled={isLoading}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? '保存中...' : isEditing ? '更新' : '作成'}
        </button>
      </div>
    </form>
  );
}
