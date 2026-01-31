'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { MemberCostResponse } from '@/schemas/budget-schema';
import { formatCurrency } from '@/types/budget';

interface CostBreakdownProps {
  memberCosts: MemberCostResponse[];
  currency?: string;
}

export function CostBreakdown({ memberCosts, currency = 'JPY' }: CostBreakdownProps) {
  const totalCost = memberCosts.reduce((sum, m) => sum + m.cost, 0);
  const totalHours = memberCosts.reduce((sum, m) => sum + m.hours, 0);

  // Sort by cost descending
  const sortedCosts = [...memberCosts].sort((a, b) => b.cost - a.cost);

  return (
    <Card>
      <CardHeader>
        <CardTitle>メンバー別コスト内訳</CardTitle>
      </CardHeader>
      <CardContent>
        {memberCosts.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            工数記録がありません
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-2 font-medium">メンバー</th>
                    <th className="text-right py-3 px-2 font-medium">工数</th>
                    <th className="text-right py-3 px-2 font-medium">時給</th>
                    <th className="text-right py-3 px-2 font-medium">コスト</th>
                    <th className="py-3 px-2 font-medium w-[150px]">割合</th>
                  </tr>
                </thead>
                <tbody>
                  {sortedCosts.map((member) => (
                    <tr key={member.member_id} className="border-b last:border-0">
                      <td className="py-3 px-2 font-medium">{member.member_name}</td>
                      <td className="py-3 px-2 text-right">{member.hours.toFixed(1)}h</td>
                      <td className="py-3 px-2 text-right">
                        {formatCurrency(member.hourly_rate, currency)}
                      </td>
                      <td className="py-3 px-2 text-right font-semibold">
                        {formatCurrency(member.cost, currency)}
                      </td>
                      <td className="py-3 px-2">
                        <div className="flex items-center gap-2">
                          <div className="flex-1 h-2 bg-gray-100 rounded-full overflow-hidden">
                            <div
                              className="h-full bg-primary rounded-full"
                              style={{ width: `${member.percentage}%` }}
                            />
                          </div>
                          <span className="text-sm text-muted-foreground w-12 text-right">
                            {member.percentage.toFixed(1)}%
                          </span>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Summary footer */}
            <div className="flex justify-between items-center mt-4 pt-4 border-t">
              <div className="text-sm text-muted-foreground">
                合計: {memberCosts.length}名
              </div>
              <div className="flex gap-6">
                <div className="text-sm">
                  <span className="text-muted-foreground">総工数: </span>
                  <span className="font-semibold">{totalHours.toFixed(1)}h</span>
                </div>
                <div className="text-sm">
                  <span className="text-muted-foreground">総コスト: </span>
                  <span className="font-semibold">{formatCurrency(totalCost, currency)}</span>
                </div>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
