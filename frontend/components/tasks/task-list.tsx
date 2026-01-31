'use client';

import { TaskResponse } from '@/schemas/task-schema';
import { TaskItem } from './task-item';
import Loading from '@/components/common/loading';

interface TaskListProps {
  tasks: TaskResponse[];
  isLoading?: boolean;
  onEdit?: (task: TaskResponse) => void;
  onDelete?: (taskId: string) => void;
  onUpdateHours?: (taskId: string, field: 'planned_hours' | 'actual_hours', value: number) => void;
}

export function TaskList({
  tasks,
  isLoading,
  onEdit,
  onDelete,
  onUpdateHours,
}: TaskListProps) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <Loading />
      </div>
    );
  }

  if (tasks.length === 0) {
    return (
      <div className="text-center py-12">
        <svg
          className="mx-auto h-12 w-12 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
          />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">タスクがありません</h3>
        <p className="mt-1 text-sm text-gray-500">
          新しいタスクを作成して工数管理を始めましょう。
        </p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th
              scope="col"
              className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              タスク名
            </th>
            <th
              scope="col"
              className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              ステータス
            </th>
            <th
              scope="col"
              className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              担当者
            </th>
            <th
              scope="col"
              className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              予定工数
            </th>
            <th
              scope="col"
              className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              実績工数
            </th>
            <th
              scope="col"
              className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              差異
            </th>
            <th scope="col" className="px-4 py-3 text-right">
              <span className="sr-only">アクション</span>
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {tasks.map((task) => (
            <TaskItem
              key={task.id}
              task={task}
              onEdit={onEdit}
              onDelete={onDelete}
              onUpdateHours={onUpdateHours}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
}
