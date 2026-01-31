'use client';

import { useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { useProject, useUpdateProject, useDeleteProject } from '@/hooks/use-projects';
import { ProjectForm } from '@/components/projects/project-form';
import { UpdateProjectFormData } from '@/schemas/project-schema';
import { projectStatusLabels, projectStatusColors } from '@/types/project';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import Loading from '@/components/common/loading';
import {
  ArrowLeft,
  Calendar,
  DollarSign,
  ListTodo,
  Clock,
  TrendingUp,
  Pencil,
  Trash2,
  ChartBar,
} from 'lucide-react';

export default function ProjectDetailPage() {
  const router = useRouter();
  const params = useParams();
  const projectId = params.id as string;

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const { data: project, isLoading, error } = useProject(projectId);
  const updateMutation = useUpdateProject();
  const deleteMutation = useDeleteProject();

  const handleUpdate = async (data: UpdateProjectFormData) => {
    try {
      await updateMutation.mutateAsync({ projectId, data });
      setIsEditDialogOpen(false);
    } catch (error) {
      console.error('プロジェクト更新に失敗しました:', error);
    }
  };

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(projectId);
      router.push('/projects');
    } catch (error) {
      console.error('プロジェクト削除に失敗しました:', error);
    }
  };

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '未設定';
    return new Date(dateStr).toLocaleDateString('ja-JP');
  };

  const formatCurrency = (amount?: number) => {
    if (amount === undefined || amount === null) return '未設定';
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loading />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px]">
        <p className="text-red-500 mb-4">プロジェクトの取得に失敗しました</p>
        <Link href="/projects">
          <Button>プロジェクト一覧に戻る</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/projects">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <div className="flex items-center space-x-3">
              <h1 className="text-2xl font-bold">{project.name}</h1>
              <Badge className={projectStatusColors[project.status]}>
                {projectStatusLabels[project.status]}
              </Badge>
            </div>
            {project.description && (
              <p className="text-muted-foreground mt-1">{project.description}</p>
            )}
          </div>
        </div>
        <div className="flex space-x-2">
          <Button variant="outline" onClick={() => setIsEditDialogOpen(true)}>
            <Pencil className="mr-2 h-4 w-4" />
            編集
          </Button>
          <Button variant="destructive" onClick={() => setIsDeleteDialogOpen(true)}>
            <Trash2 className="mr-2 h-4 w-4" />
            削除
          </Button>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">タスク数</CardTitle>
            <ListTodo className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{project.stats?.total_tasks ?? 0}</div>
            <p className="text-xs text-muted-foreground">
              {project.stats?.completed_tasks ?? 0} 件完了
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">予定工数</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {project.stats?.total_planned_hours?.toFixed(1) ?? 0}h
            </div>
            <p className="text-xs text-muted-foreground">
              実績: {project.stats?.total_actual_hours?.toFixed(1) ?? 0}h
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">進捗率</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {(project.stats?.completion_rate ?? 0).toFixed(0)}%
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
              <div
                className="bg-blue-600 h-2 rounded-full"
                style={{ width: `${project.stats?.completion_rate ?? 0}%` }}
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">予算</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(project.budget_amount)}</div>
            <p className="text-xs text-muted-foreground">設定済み</p>
          </CardContent>
        </Card>
      </div>

      {/* Project Info */}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>プロジェクト情報</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center space-x-3">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">期間</p>
                <p className="text-sm text-muted-foreground">
                  {formatDate(project.start_date)} 〜 {formatDate(project.end_date)}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <DollarSign className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">予算</p>
                <p className="text-sm text-muted-foreground">
                  {formatCurrency(project.budget_amount)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>クイックアクション</CardTitle>
            <CardDescription>プロジェクトの詳細な管理ページへ移動</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <Link href={`/projects/${projectId}/tasks`} className="block">
              <Button variant="outline" className="w-full justify-start">
                <ListTodo className="mr-2 h-4 w-4" />
                タスク管理
              </Button>
            </Link>
            <Link href={`/projects/${projectId}/budget`} className="block">
              <Button variant="outline" className="w-full justify-start">
                <DollarSign className="mr-2 h-4 w-4" />
                収支管理
              </Button>
            </Link>
            <Link href={`/projects/${projectId}/analytics`} className="block">
              <Button variant="outline" className="w-full justify-start">
                <ChartBar className="mr-2 h-4 w-4" />
                分析・レポート
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>プロジェクト編集</DialogTitle>
            <DialogDescription>プロジェクトの情報を編集します。</DialogDescription>
          </DialogHeader>
          <ProjectForm
            project={project}
            onSubmit={handleUpdate}
            onCancel={() => setIsEditDialogOpen(false)}
            isLoading={updateMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>プロジェクトを削除しますか？</DialogTitle>
            <DialogDescription>
              「{project.name}」を削除します。関連するタスクや収支データも削除されます。
              この操作は取り消せません。
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end space-x-4 pt-4">
            <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>
              キャンセル
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
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
