'use client';

import { useState, use } from 'react';
import { useRouter } from 'next/navigation';
import { useTasks, useProjectSummary, useCreateTask, useUpdateTask, useDeleteTask } from '@/hooks/use-tasks';
import { TaskList } from '@/components/tasks/task-list';
import { TaskForm } from '@/components/tasks/task-form';
import { PlanActualSummary } from '@/components/tasks/plan-actual-summary';
import Pagination from '@/components/common/pagination';
import Loading from '@/components/common/loading';
import { TaskResponse, CreateTaskFormData, UpdateTaskFormData, taskStatusLabels } from '@/schemas/task-schema';

interface TasksPageProps {
  params: Promise<{ id: string }>;
}

export default function TasksPage(props: TasksPageProps) {
  const params = use(props.params);
  const projectId = params.id;
  const router = useRouter();

  // State
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<TaskResponse | null>(null);

  // Queries
  const { data: tasksData, isLoading: isLoadingTasks } = useTasks(projectId, {
    page,
    per_page: 20,
    status: statusFilter || undefined,
  });

  const { data: summary, isLoading: isLoadingSummary } = useProjectSummary(projectId);

  // Mutations
  const createTask = useCreateTask(projectId);
  const updateTask = useUpdateTask(projectId);
  const deleteTask = useDeleteTask(projectId);

  // Handlers
  const handleSubmitTask = (data: CreateTaskFormData | UpdateTaskFormData) => {
    if (editingTask) {
      // Update existing task
      updateTask.mutate(
        { taskId: editingTask.id, data: data as UpdateTaskFormData },
        {
          onSuccess: () => {
            setEditingTask(null);
            setIsFormOpen(false);
          },
        }
      );
    } else {
      // Create new task
      createTask.mutate(data as CreateTaskFormData, {
        onSuccess: () => {
          setIsFormOpen(false);
        },
      });
    }
  };

  const handleDeleteTask = (taskId: string) => {
    if (confirm('このタスクを削除しますか？')) {
      deleteTask.mutate(taskId);
    }
  };

  const handleEditTask = (task: TaskResponse) => {
    setEditingTask(task);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setIsFormOpen(false);
    setEditingTask(null);
  };

  if (isLoadingTasks && isLoadingSummary) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loading />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* ヘッダー */}
      <div className="flex items-center justify-between">
        <div>
          <button
            onClick={() => router.back()}
            className="text-sm text-gray-600 hover:text-gray-900 mb-2 flex items-center"
          >
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            戻る
          </button>
          <h1 className="text-2xl font-bold text-gray-900">タスク管理</h1>
          <p className="text-sm text-gray-500 mt-1">予定工数と実績工数を管理します</p>
        </div>
        <button
          onClick={() => setIsFormOpen(true)}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          新規タスク
        </button>
      </div>

      {/* 予実サマリー */}
      {summary && <PlanActualSummary summary={summary} />}

      {/* フィルター */}
      <div className="bg-white rounded-lg shadow px-6 py-4">
        <div className="flex items-center gap-4">
          <label htmlFor="status-filter" className="text-sm font-medium text-gray-700">
            ステータス:
          </label>
          <select
            id="status-filter"
            value={statusFilter}
            onChange={(e) => {
              setStatusFilter(e.target.value);
              setPage(1);
            }}
            className="rounded-md border border-gray-300 px-3 py-1.5 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="">すべて</option>
            {Object.entries(taskStatusLabels).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* タスクリスト */}
      <div className="bg-white rounded-lg shadow">
        <TaskList
          tasks={tasksData?.tasks || []}
          isLoading={isLoadingTasks}
          onEdit={handleEditTask}
          onDelete={handleDeleteTask}
        />

        {/* ページネーション */}
        {tasksData && tasksData.pagination.total_pages > 1 && (
          <div className="px-6 py-4 border-t border-gray-200">
            <Pagination
              currentPage={page}
              totalPages={tasksData.pagination.total_pages}
              onPageChange={setPage}
            />
          </div>
        )}
      </div>

      {/* タスクフォームモーダル */}
      {isFormOpen && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
            {/* 背景 */}
            <div
              className="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75"
              onClick={handleCloseForm}
            />

            {/* モーダル */}
            <div className="inline-block w-full max-w-lg p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-white shadow-xl rounded-lg">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold text-gray-900">
                  {editingTask ? 'タスクを編集' : '新規タスク'}
                </h2>
                <button
                  onClick={handleCloseForm}
                  className="text-gray-400 hover:text-gray-500"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <TaskForm
                task={editingTask || undefined}
                onSubmit={handleSubmitTask}
                onCancel={handleCloseForm}
                isLoading={createTask.isPending || updateTask.isPending}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
