'use client';

import { useParams, useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { BudgetSummary } from '@/components/budget/budget-summary';
import { CostBreakdown } from '@/components/budget/cost-breakdown';
import { RevenueForm } from '@/components/budget/revenue-form';
import Loading from '@/components/common/loading';
import { useBudgetSummary } from '@/hooks/use-budget';
import { ArrowLeft } from 'lucide-react';

export default function BudgetPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const { data: budgetSummary, isLoading, error } = useBudgetSummary(projectId);

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return (
      <div className="p-8">
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold text-red-600">エラーが発生しました</h2>
          <p className="text-muted-foreground mt-2">{(error as Error).message}</p>
          <Button
            variant="outline"
            className="mt-4"
            onClick={() => router.push(`/projects/${projectId}`)}
          >
            プロジェクトに戻る
          </Button>
        </div>
      </div>
    );
  }

  if (!budgetSummary) {
    return (
      <div className="p-8">
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold">収支データが見つかりません</h2>
          <Button
            variant="outline"
            className="mt-4"
            onClick={() => router.push(`/projects/${projectId}`)}
          >
            プロジェクトに戻る
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => router.push(`/projects/${projectId}`)}
          >
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-2xl font-bold">収支管理</h1>
            <p className="text-muted-foreground">{budgetSummary.project_name}</p>
          </div>
        </div>
      </div>

      {/* Revenue Form */}
      <RevenueForm
        projectId={projectId}
        currentRevenue={budgetSummary.budget.revenue}
        currency={budgetSummary.budget.currency}
      />

      {/* Budget Summary */}
      <BudgetSummary budget={budgetSummary} />

      {/* Cost Breakdown by Member */}
      <CostBreakdown
        memberCosts={budgetSummary.member_costs}
        currency={budgetSummary.budget.currency}
      />
    </div>
  );
}
