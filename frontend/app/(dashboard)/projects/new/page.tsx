'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useCreateProject } from '@/hooks/use-projects';
import { ProjectForm } from '@/components/projects/project-form';
import { CreateProjectFormData, UpdateProjectFormData } from '@/schemas/project-schema';
import { CreateProjectRequest } from '@/types/project';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';

export default function NewProjectPage() {
  const router = useRouter();
  const createMutation = useCreateProject();

  const handleSubmit = async (data: CreateProjectFormData | UpdateProjectFormData) => {
    try {
      const project = await createMutation.mutateAsync(data as CreateProjectRequest);
      router.push(`/projects/${project.id}`);
    } catch (error) {
      console.error('プロジェクト作成に失敗しました:', error);
    }
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div className="flex items-center space-x-4">
        <Link href="/projects">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-2xl font-bold">新規プロジェクト作成</h1>
          <p className="text-muted-foreground">
            新しいプロジェクトを作成します
          </p>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>プロジェクト情報</CardTitle>
        </CardHeader>
        <CardContent>
          <ProjectForm
            onSubmit={handleSubmit}
            onCancel={() => router.push('/projects')}
            isLoading={createMutation.isPending}
          />
        </CardContent>
      </Card>
    </div>
  );
}
