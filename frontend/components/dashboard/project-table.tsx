'use client';

import { useState } from 'react';
import Link from 'next/link';
import { ArrowDown, ArrowUp, ArrowUpDown, Eye, Pencil } from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import Pagination from '@/components/common/pagination';
import Loading from '@/components/common/loading';
import { useProjects } from '@/hooks/use-projects';
import {
  Project,
  ProjectStatus,
  projectStatusColors,
  projectStatusLabels,
} from '@/types/project';

type SortKey = 'name' | 'status' | 'start_date' | 'created_at' | 'profit_rate';

const sortableColumns: { key: SortKey; label: string }[] = [
  { key: 'name', label: 'プロジェクト名' },
  { key: 'status', label: 'ステータス' },
  { key: 'start_date', label: '期間' },
  { key: 'profit_rate', label: '利益率' },
];

const PER_PAGE = 10;

function formatPeriod(project: Project): string {
  if (!project.start_date && !project.end_date) return '-';
  return `${project.start_date ?? ''} 〜 ${project.end_date ?? ''}`;
}

function formatBudget(project: Project): string {
  if (project.budget_amount == null) return '-';
  return `¥${project.budget_amount.toLocaleString()}`;
}

export function ProjectTable() {
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<'all' | ProjectStatus>('all');
  const [sort, setSort] = useState<SortKey>('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');

  const { data, isLoading, isError, error } = useProjects({
    page,
    per_page: PER_PAGE,
    status: status === 'all' ? undefined : status,
    sort,
    order,
  });

  const toggleSort = (key: SortKey) => {
    if (sort === key) {
      setOrder(order === 'asc' ? 'desc' : 'asc');
    } else {
      setSort(key);
      setOrder(key === 'name' ? 'asc' : 'desc');
    }
    setPage(1);
  };

  const sortIcon = (key: SortKey) => {
    if (sort !== key) return <ArrowUpDown className="ml-1 h-3 w-3" />;
    return order === 'asc' ? (
      <ArrowUp className="ml-1 h-3 w-3" />
    ) : (
      <ArrowDown className="ml-1 h-3 w-3" />
    );
  };

  if (isError) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-red-800">
        プロジェクトの取得に失敗しました: {error instanceof Error ? error.message : ''}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* ステータスフィルタ */}
      <div className="flex items-center justify-between">
        <Select
          value={status}
          onValueChange={(value) => {
            setStatus(value as 'all' | ProjectStatus);
            setPage(1);
          }}
        >
          <SelectTrigger className="w-44">
            <SelectValue placeholder="ステータスで絞り込み" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">すべて</SelectItem>
            {(Object.keys(projectStatusLabels) as ProjectStatus[]).map((value) => (
              <SelectItem key={value} value={value}>
                {projectStatusLabels[value]}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {data && (
          <span className="text-sm text-muted-foreground">全{data.pagination.total}件</span>
        )}
      </div>

      {isLoading ? (
        <Loading />
      ) : (
        <>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  {sortableColumns.map(({ key, label }) => (
                    <TableHead key={key}>
                      <button
                        type="button"
                        className="flex items-center font-medium hover:text-foreground"
                        onClick={() => toggleSort(key)}
                      >
                        {label}
                        {sortIcon(key)}
                      </button>
                    </TableHead>
                  ))}
                  <TableHead>予算</TableHead>
                  <TableHead className="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data?.projects.length ? (
                  data.projects.map((project) => (
                    <TableRow key={project.id}>
                      <TableCell className="font-medium">
                        <Link href={`/projects/${project.id}`} className="hover:underline">
                          {project.name}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <Badge className={projectStatusColors[project.status]}>
                          {projectStatusLabels[project.status]}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {formatPeriod(project)}
                      </TableCell>
                      <TableCell className="text-muted-foreground">-</TableCell>
                      <TableCell>{formatBudget(project)}</TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-1">
                          <Button variant="ghost" size="sm" asChild>
                            <Link href={`/projects/${project.id}`}>
                              <Eye className="h-4 w-4" />
                              <span className="sr-only">詳細</span>
                            </Link>
                          </Button>
                          <Button variant="ghost" size="sm" asChild>
                            <Link href={`/projects/${project.id}?edit=1`}>
                              <Pencil className="h-4 w-4" />
                              <span className="sr-only">編集</span>
                            </Link>
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                      プロジェクトがありません
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          {data && data.pagination.total_pages > 1 && (
            <Pagination
              currentPage={page}
              totalPages={data.pagination.total_pages}
              onPageChange={setPage}
            />
          )}
        </>
      )}
    </div>
  );
}
