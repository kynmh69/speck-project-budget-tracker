'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';
import { useProjects, useCreateProject, useUpdateProject, useDeleteProject } from '@/hooks/use-projects';
import { Project, CreateProjectRequest, UpdateProjectRequest } from '@/types/project';
import { CreateProjectFormData, UpdateProjectFormData } from '@/schemas/project-schema';
import { ProjectCard } from '@/components/projects/project-card';
import { ProjectForm } from '@/components/projects/project-form';
import { ProjectFilters, ProjectFilters as ProjectFiltersType } from '@/components/projects/project-filters';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import Loading from '@/components/common/loading';
import Pagination from '@/components/common/pagination';

export default function ProjectsPage() {
  const router = useRouter();
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<ProjectFiltersType>({});
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [deletingProject, setDeletingProject] = useState<Project | null>(null);

  const perPage = 12;

  const { data, isLoading, error } = useProjects({
    page,
    per_page: perPage,
    ...filters,
  });

  const createMutation = useCreateProject();
  const updateMutation = useUpdateProject();
  const deleteMutation = useDeleteProject();

  const handleFilterChange = (newFilters: ProjectFiltersType) => {
    setFilters(newFilters);
    setPage(1);
  };

  const handleCreateProject = async (formData: CreateProjectFormData | UpdateProjectFormData) => {
    try {
      await createMutation.mutateAsync(formData as CreateProjectRequest);
      setIsCreateDialogOpen(false);
    } catch (error) {
      console.error('プロジェクト作成に失敗しました:', error);
    }
  };

  const handleUpdateProject = async (formData: CreateProjectFormData | UpdateProjectFormData) => {
    if (!editingProject) return;
    try {
      await updateMutation.mutateAsync({
        projectId: editingProject.id,
        data: formData as UpdateProjectRequest,
      });
      setEditingProject(null);
    } catch (error) {
      console.error('プロジェクト更新に失敗しました:', error);
    }
  };

  const handleDeleteProject = async () => {
    if (!deletingProject) return;
    try {
      await deleteMutation.mutateAsync(deletingProject.id);
      setDeletingProject(null);
    } catch (error) {
      console.error('プロジェクト削除に失敗しました:', error);
    }
  };

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <p className="text-red-500 mb-4">プロジェクトの取得に失敗しました</p>
          <Button onClick={() => window.location.reload()}>再読み込み</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">プロジェクト</h1>
          <p className="text-muted-foreground">プロジェクトの管理と進捗状況を確認できます</p>
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          新規プロジェクト
        </Button>
      </div>

      {/* Filters */}
      <ProjectFilters onFilterChange={handleFilterChange} initialFilters={filters} />

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center min-h-[400px]">
          <Loading />
        </div>
      ) : data?.projects.length === 0 ? (
        <div className="flex flex-col items-center justify-center min-h-[400px] text-center">
          <p className="text-muted-foreground mb-4">
            {Object.keys(filters).some((k) => filters[k as keyof ProjectFiltersType])
              ? '条件に一致するプロジェクトが見つかりません'
              : 'プロジェクトがありません'}
          </p>
          <Button onClick={() => setIsCreateDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            最初のプロジェクトを作成
          </Button>
        </div>
      ) : (
        <>
          {/* Project Grid */}
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {data?.projects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onEdit={setEditingProject}
                onDelete={setDeletingProject}
              />
            ))}
          </div>

          {/* Pagination */}
          {data && data.pagination.total_pages > 1 && (
            <Pagination
              currentPage={page}
              totalPages={data.pagination.total_pages}
              onPageChange={setPage}
            />
          )}
        </>
      )}

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>新規プロジェクト作成</DialogTitle>
            <DialogDescription>
              新しいプロジェクトを作成します。必要な情報を入力してください。
            </DialogDescription>
          </DialogHeader>
          <ProjectForm
            onSubmit={handleCreateProject}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={!!editingProject} onOpenChange={(open: boolean) => !open && setEditingProject(null)}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>プロジェクト編集</DialogTitle>
            <DialogDescription>
              プロジェクトの情報を編集します。
            </DialogDescription>
          </DialogHeader>
          {editingProject && (
            <ProjectForm
              project={editingProject}
              onSubmit={handleUpdateProject}
              onCancel={() => setEditingProject(null)}
              isLoading={updateMutation.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!deletingProject} onOpenChange={(open: boolean) => !open && setDeletingProject(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>プロジェクトを削除しますか？</DialogTitle>
            <DialogDescription>
              「{deletingProject?.name}」を削除します。この操作は取り消せません。
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end space-x-4 pt-4">
            <Button variant="outline" onClick={() => setDeletingProject(null)}>
              キャンセル
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteProject}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? '削除中...' : '削除'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
