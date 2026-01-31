'use client';

import { ProjectSummary } from '@/schemas/task-schema';
import { cn } from '@/lib/utils';

interface PlanActualSummaryProps {
  summary: ProjectSummary;
  className?: string;
}

export function PlanActualSummary({ summary, className }: PlanActualSummaryProps) {
  const isOverBudget = summary.is_over_budget;

  return (
    <div className={cn('bg-white rounded-lg shadow', className)}>
      <div className="px-6 py-4 border-b border-gray-200">
        <h3 className="text-lg font-semibold text-gray-900">予実サマリー</h3>
      </div>

      <div className="p-6">
        {/* 工数サマリー */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          {/* 予定工数 */}
          <div className="bg-blue-50 rounded-lg p-4">
            <p className="text-sm font-medium text-blue-600">予定工数</p>
            <p className="text-2xl font-bold text-blue-900">
              {summary.total_planned_hours.toFixed(1)}h
            </p>
          </div>

          {/* 実績工数 */}
          <div className="bg-green-50 rounded-lg p-4">
            <p className="text-sm font-medium text-green-600">実績工数</p>
            <p className="text-2xl font-bold text-green-900">
              {summary.total_actual_hours.toFixed(1)}h
            </p>
          </div>

          {/* 差異 */}
          <div className={cn(
            'rounded-lg p-4',
            isOverBudget ? 'bg-red-50' : 'bg-emerald-50'
          )}>
            <p className={cn(
              'text-sm font-medium',
              isOverBudget ? 'text-red-600' : 'text-emerald-600'
            )}>
              差異
            </p>
            <p className={cn(
              'text-2xl font-bold',
              isOverBudget ? 'text-red-900' : 'text-emerald-900'
            )}>
              {summary.variance_hours > 0 ? '+' : ''}
              {summary.variance_hours.toFixed(1)}h
              <span className="text-sm font-normal ml-2">
                ({summary.variance_percentage > 0 ? '+' : ''}{summary.variance_percentage.toFixed(1)}%)
              </span>
            </p>
          </div>
        </div>

        {/* 警告メッセージ */}
        {isOverBudget && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <div className="flex items-center">
              <svg
                className="h-5 w-5 text-red-400 mr-2"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                  clipRule="evenodd"
                />
              </svg>
              <p className="text-sm font-medium text-red-800">
                工数超過: 予定工数を{summary.variance_hours.toFixed(1)}時間（{summary.variance_percentage.toFixed(1)}%）超過しています
              </p>
            </div>
          </div>
        )}

        {/* タスク進捗 */}
        <div className="border-t border-gray-200 pt-6">
          <h4 className="text-sm font-medium text-gray-700 mb-4">タスク進捗</h4>
          
          {/* プログレスバー */}
          <div className="mb-4">
            <div className="flex justify-between text-sm text-gray-600 mb-1">
              <span>完了率</span>
              <span>{summary.completion_rate.toFixed(1)}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2.5">
              <div
                className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
                style={{ width: `${Math.min(summary.completion_rate, 100)}%` }}
              />
            </div>
          </div>

          {/* タスク統計 */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="text-center">
              <p className="text-2xl font-bold text-gray-900">{summary.total_tasks}</p>
              <p className="text-xs text-gray-500">合計タスク</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold text-green-600">{summary.completed_tasks}</p>
              <p className="text-xs text-gray-500">完了</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold text-blue-600">{summary.in_progress_tasks}</p>
              <p className="text-xs text-gray-500">進行中</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold text-gray-600">{summary.todo_tasks}</p>
              <p className="text-xs text-gray-500">未着手</p>
            </div>
          </div>

          {summary.blocked_tasks > 0 && (
            <div className="mt-4 text-center">
              <p className="text-lg font-bold text-red-600">{summary.blocked_tasks}</p>
              <p className="text-xs text-gray-500">ブロック中</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
