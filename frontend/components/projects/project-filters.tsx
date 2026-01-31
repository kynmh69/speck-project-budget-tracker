'use client';

import { useState } from 'react';
import { ProjectStatus, projectStatusLabels } from '@/types/project';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Search, X } from 'lucide-react';

interface ProjectFiltersProps {
  onFilterChange: (filters: ProjectFilters) => void;
  initialFilters?: ProjectFilters;
}

export interface ProjectFilters {
  search?: string;
  status?: string;
  sort?: string;
  order?: string;
}

export function ProjectFilters({ onFilterChange, initialFilters }: ProjectFiltersProps) {
  const [search, setSearch] = useState(initialFilters?.search || '');
  const [status, setStatus] = useState(initialFilters?.status || '');
  const [sort, setSort] = useState(initialFilters?.sort || 'created_at');
  const [order, setOrder] = useState(initialFilters?.order || 'desc');

  const handleSearchChange = (value: string) => {
    setSearch(value);
    onFilterChange({ search: value, status, sort, order });
  };

  const handleStatusChange = (value: string) => {
    const newStatus = value === 'all' ? '' : value;
    setStatus(newStatus);
    onFilterChange({ search, status: newStatus, sort, order });
  };

  const handleSortChange = (value: string) => {
    setSort(value);
    onFilterChange({ search, status, sort: value, order });
  };

  const handleOrderChange = (value: string) => {
    setOrder(value);
    onFilterChange({ search, status, sort, order: value });
  };

  const clearFilters = () => {
    setSearch('');
    setStatus('');
    setSort('created_at');
    setOrder('desc');
    onFilterChange({ search: '', status: '', sort: 'created_at', order: 'desc' });
  };

  const hasActiveFilters = search || status || sort !== 'created_at' || order !== 'desc';

  const statusOptions: { value: string; label: string }[] = [
    { value: 'all', label: 'すべて' },
    { value: 'planning', label: projectStatusLabels.planning },
    { value: 'in_progress', label: projectStatusLabels.in_progress },
    { value: 'completed', label: projectStatusLabels.completed },
    { value: 'on_hold', label: projectStatusLabels.on_hold },
  ];

  const sortOptions = [
    { value: 'created_at', label: '作成日' },
    { value: 'name', label: '名前' },
    { value: 'start_date', label: '開始日' },
    { value: 'end_date', label: '終了日' },
    { value: 'status', label: 'ステータス' },
  ];

  const orderOptions = [
    { value: 'desc', label: '降順' },
    { value: 'asc', label: '昇順' },
  ];

  return (
    <div className="flex flex-col md:flex-row gap-4 p-4 bg-white rounded-lg border">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="プロジェクト名で検索..."
          value={search}
          onChange={(e) => handleSearchChange(e.target.value)}
          className="pl-10"
        />
      </div>

      <Select value={status || 'all'} onValueChange={handleStatusChange}>
        <SelectTrigger className="w-full md:w-40">
          <SelectValue placeholder="ステータス" />
        </SelectTrigger>
        <SelectContent>
          {statusOptions.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={sort} onValueChange={handleSortChange}>
        <SelectTrigger className="w-full md:w-32">
          <SelectValue placeholder="並び替え" />
        </SelectTrigger>
        <SelectContent>
          {sortOptions.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select value={order} onValueChange={handleOrderChange}>
        <SelectTrigger className="w-full md:w-24">
          <SelectValue placeholder="順序" />
        </SelectTrigger>
        <SelectContent>
          {orderOptions.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      {hasActiveFilters && (
        <Button variant="ghost" size="icon" onClick={clearFilters}>
          <X className="h-4 w-4" />
        </Button>
      )}
    </div>
  );
}
