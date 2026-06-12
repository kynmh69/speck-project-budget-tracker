'use client';

import {
  Briefcase,
  CheckCircle2,
  DollarSign,
  Percent,
  PlayCircle,
  TrendingUp,
} from 'lucide-react';

import Loading from '@/components/common/loading';
import { KpiCard } from '@/components/dashboard/kpi-card';
import { ProjectTable } from '@/components/dashboard/project-table';
import { useCurrentUser } from '@/hooks/use-auth';
import { useDashboard } from '@/hooks/use-dashboard';
import { formatCurrency } from '@/types/budget';

export default function DashboardPage() {
  const { data: user, isLoading: isUserLoading } = useCurrentUser();
  const { data: dashboard, isLoading: isDashboardLoading, isError, error } = useDashboard();

  if (isUserLoading || isDashboardLoading) {
    return <Loading />;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">ダッシュボード</h1>
        {user && <p className="text-muted-foreground mt-1">ようこそ、{user.name}さん</p>}
      </div>

      {isError && (
        <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-red-800">
          ダッシュボードの取得に失敗しました: {error instanceof Error ? error.message : ''}
        </div>
      )}

      {dashboard && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
          <KpiCard title="総プロジェクト数" value={dashboard.total_projects} icon={Briefcase} />
          <KpiCard title="進行中" value={dashboard.active_projects} icon={PlayCircle} />
          <KpiCard title="完了" value={dashboard.completed_projects} icon={CheckCircle2} />
          <KpiCard
            title="総売上"
            value={formatCurrency(dashboard.total_revenue)}
            icon={DollarSign}
          />
          <KpiCard
            title="総利益"
            value={formatCurrency(dashboard.total_profit)}
            icon={TrendingUp}
            valueClassName={dashboard.total_profit < 0 ? 'text-red-600' : 'text-green-600'}
          />
          <KpiCard
            title="平均利益率"
            value={`${dashboard.average_profit_rate.toFixed(1)}%`}
            icon={Percent}
          />
        </div>
      )}

      <section className="space-y-3">
        <h2 className="text-xl font-semibold">プロジェクト一覧</h2>
        <ProjectTable />
      </section>
    </div>
  );
}
