'use client';

import Link from 'next/link';
import { Project, projectStatusLabels, projectStatusColors } from '@/types/project';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Calendar, DollarSign, MoreHorizontal, Pencil, Trash2 } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

interface ProjectCardProps {
  project: Project;
  onEdit?: (project: Project) => void;
  onDelete?: (project: Project) => void;
}

export function ProjectCard({ project, onEdit, onDelete }: ProjectCardProps) {
  const formatDate = (dateStr?: string) => {
    if (!dateStr) return null;
    return new Date(dateStr).toLocaleDateString('ja-JP');
  };

  const formatCurrency = (amount?: number) => {
    if (amount === undefined || amount === null) return null;
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
        <div className="space-y-1">
          <Link href={`/projects/${project.id}`}>
            <CardTitle className="text-lg font-semibold hover:text-blue-600 cursor-pointer">
              {project.name}
            </CardTitle>
          </Link>
          <Badge className={projectStatusColors[project.status]}>
            {projectStatusLabels[project.status]}
          </Badge>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onEdit?.(project)}>
              <Pencil className="mr-2 h-4 w-4" />
              編集
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => onDelete?.(project)}
              className="text-red-600"
            >
              <Trash2 className="mr-2 h-4 w-4" />
              削除
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardHeader>
      <CardContent>
        {project.description && (
          <p className="text-sm text-muted-foreground line-clamp-2 mb-4">
            {project.description}
          </p>
        )}
        <div className="space-y-2 text-sm">
          {(project.start_date || project.end_date) && (
            <div className="flex items-center text-muted-foreground">
              <Calendar className="mr-2 h-4 w-4" />
              <span>
                {formatDate(project.start_date) || '未定'}
                {' 〜 '}
                {formatDate(project.end_date) || '未定'}
              </span>
            </div>
          )}
          {project.budget_amount !== undefined && project.budget_amount !== null && (
            <div className="flex items-center text-muted-foreground">
              <DollarSign className="mr-2 h-4 w-4" />
              <span>予算: {formatCurrency(project.budget_amount)}</span>
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="pt-0">
        <Link href={`/projects/${project.id}`} className="w-full">
          <Button variant="outline" className="w-full">
            詳細を見る
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
