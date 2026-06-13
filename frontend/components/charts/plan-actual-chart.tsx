'use client';

import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Legend,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

import { ChartContainer } from '@/components/charts/chart-container';
import { chartColors, chartDefaults, formatHours } from '@/components/charts/chart-theme';
import { PlanActualTaskData } from '@/types/analytics';

interface PlanActualChartProps {
  data?: PlanActualTaskData[];
  isLoading?: boolean;
  errorMessage?: string;
}

/**
 * タスク別の予定工数と実績工数を比較する棒グラフ (T180)
 * 実績が予定を超過したタスクは実績バーを警告色で表示する
 */
export function PlanActualChart({ data = [], isLoading, errorMessage }: PlanActualChartProps) {
  return (
    <ChartContainer
      title="予実比較（タスク別）"
      isLoading={isLoading}
      errorMessage={errorMessage}
      isEmpty={data.length === 0}
      emptyMessage="タスクがありません"
    >
      <BarChart data={data} margin={chartDefaults.margin}>
        <CartesianGrid strokeDasharray="3 3" stroke={chartDefaults.gridStroke} />
        <XAxis dataKey="task_name" fontSize={chartDefaults.axisFontSize} />
        <YAxis fontSize={chartDefaults.axisFontSize} tickFormatter={formatHours} />
        <Tooltip formatter={(value: number) => formatHours(value)} />
        <Legend />
        <Bar dataKey="planned_hours" name="予定工数" fill={chartColors.planned} />
        <Bar dataKey="actual_hours" name="実績工数" fill={chartColors.actual}>
          {data.map((entry) => (
            <Cell
              key={entry.task_name}
              fill={entry.variance > 0 ? chartColors.deficit : chartColors.actual}
            />
          ))}
        </Bar>
      </BarChart>
    </ChartContainer>
  );
}
