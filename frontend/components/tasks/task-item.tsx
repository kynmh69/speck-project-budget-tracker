'use client';

import { TaskResponse, taskStatusLabels, taskStatusColors } from '@/schemas/task-schema';
import { cn } from '@/lib/utils';

interface TaskItemProps {
  task: TaskResponse;
  onEdit?: (task: TaskResponse) => void;
  onDelete?: (taskId: string) => void;
  onUpdateHours?: (taskId: string, field: 'planned_hours' | 'actual_hours', value: number) => void;
}

export function TaskItem({ task, onEdit, onDelete, onUpdateHours }: TaskItemProps) {
  const isOverBudget = task.variance_hours > 0;
  const varianceColor = isOverBudget ? 'text-red-600' : 'text-green-600';

  return (
    <tr className="hover:bg-gray-50">
      {/* タスク名 */}
      <td className="px-4 py-3">
        <div>
          <p className="font-medium text-gray-900">{task.name}</p>
          {task.description && (
            <p className="text-sm text-gray-500 truncate max-w-xs">{task.description}</p>
          )}
        </div>
      </td>

      {/* ステータス */}
      <td className="px-4 py-3">
        <span
          className={cn(
            'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
            taskStatusColors[task.status]
          )}
        >
          {taskStatusLabels[task.status]}
        </span>
      </td>

      {/* 担当者 */}
      <td className="px-4 py-3 text-sm text-gray-500">
        {task.assignee?.name || '-'}
      </td>

      {/* 予定工数 */}
      <td className="px-4 py-3 text-sm text-gray-900 text-right">
        {task.planned_hours.toFixed(1)}h
      </td>

      {/* 実績工数 */}
      <td className="px-4 py-3 text-sm text-gray-900 text-right">
        {task.actual_hours.toFixed(1)}h
      </td>

      {/* 差異 */}
      <td className={cn('px-4 py-3 text-sm text-right font-medium', varianceColor)}>
        {task.variance_hours > 0 ? '+' : ''}
        {task.variance_hours.toFixed(1)}h
        {task.variance_percentage !== 0 && (
          <span className="ml-1 text-xs">
            ({task.variance_percentage > 0 ? '+' : ''}{task.variance_percentage.toFixed(1)}%)
          </span>
        )}
      </td>

      {/* アクション */}
      <td className="px-4 py-3 text-right">
        <div className="flex justify-end gap-2">
          {onEdit && (
            <button
              onClick={() => onEdit(task)}
              className="text-blue-600 hover:text-blue-800 text-sm font-medium"
            >
              編集
            </button>
          )}
          {onDelete && (
            <button
              onClick={() => onDelete(task.id)}
              className="text-red-600 hover:text-red-800 text-sm font-medium"
            >
              削除
            </button>
          )}
        </div>
      </td>
    </tr>
  );
}
