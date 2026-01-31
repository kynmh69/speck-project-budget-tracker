'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { BudgetSummaryResponse } from '@/schemas/budget-schema';
import { formatCurrency, getProfitStatusColor, getProfitBadgeColor } from '@/types/budget';
import { TrendingUp, TrendingDown, DollarSign, Clock, Calculator } from 'lucide-react';

interface BudgetSummaryProps {
  budget: BudgetSummaryResponse;
}

export function BudgetSummary({ budget }: BudgetSummaryProps) {
  const { budget: budgetData, cost_breakdown, warning_message } = budget;

  return (
    <div className="space-y-6">
      {/* Warning message */}
      {warning_message && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 flex items-center gap-3">
          <TrendingDown className="h-5 w-5 text-red-600" />
          <span className="text-red-800">{warning_message}</span>
        </div>
      )}

      {/* Main budget cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Revenue */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">売上</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatCurrency(budgetData.revenue, budgetData.currency)}
            </div>
          </CardContent>
        </Card>

        {/* Total Cost */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">総コスト</CardTitle>
            <Calculator className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatCurrency(budgetData.total_cost, budgetData.currency)}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              {cost_breakdown.total_hours.toFixed(1)}時間 × 平均{formatCurrency(cost_breakdown.average_rate, budgetData.currency)}/時
            </p>
          </CardContent>
        </Card>

        {/* Profit */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">利益</CardTitle>
            {budgetData.profit >= 0 ? (
              <TrendingUp className="h-4 w-4 text-green-600" />
            ) : (
              <TrendingDown className="h-4 w-4 text-red-600" />
            )}
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getProfitStatusColor(budgetData.profit)}`}>
              {formatCurrency(budgetData.profit, budgetData.currency)}
            </div>
            <Badge className={`mt-2 ${getProfitBadgeColor(budgetData.is_deficit)}`}>
              {budgetData.is_deficit ? '赤字' : '黒字'}
            </Badge>
          </CardContent>
        </Card>

        {/* Profit Rate */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">利益率</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getProfitStatusColor(budgetData.profit_rate)}`}>
              {budgetData.profit_rate.toFixed(1)}%
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              総工数: {cost_breakdown.total_hours.toFixed(1)}時間
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Cost breakdown summary */}
      <Card>
        <CardHeader>
          <CardTitle>コスト内訳サマリー</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex flex-col">
              <span className="text-sm text-muted-foreground">人件費</span>
              <span className="text-lg font-semibold">
                {formatCurrency(cost_breakdown.labor_cost, budgetData.currency)}
              </span>
            </div>
            <div className="flex flex-col">
              <span className="text-sm text-muted-foreground">総工数</span>
              <span className="text-lg font-semibold">{cost_breakdown.total_hours.toFixed(1)}時間</span>
            </div>
            <div className="flex flex-col">
              <span className="text-sm text-muted-foreground">平均単価</span>
              <span className="text-lg font-semibold">
                {formatCurrency(cost_breakdown.average_rate, budgetData.currency)}/時間
              </span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
