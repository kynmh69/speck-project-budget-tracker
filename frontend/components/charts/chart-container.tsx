'use client';

import { ReactNode } from 'react';
import { ResponsiveContainer } from 'recharts';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import Loading from '@/components/common/loading';
import { chartDefaults } from '@/components/charts/chart-theme';

interface ChartContainerProps {
  /** グラフのタイトル */
  title: string;
  /** データ取得中の表示制御 */
  isLoading?: boolean;
  /** エラーメッセージ（指定時はグラフの代わりに表示） */
  errorMessage?: string;
  /** データが空のときのメッセージ */
  emptyMessage?: string;
  /** データが空かどうか */
  isEmpty?: boolean;
  /** グラフの高さ(px) */
  height?: number;
  /** Rechartsのチャート要素 */
  children: ReactNode;
}

/**
 * グラフ共通のラッパー (T179)
 * Card + タイトル + ResponsiveContainer とローディング・エラー・
 * 空状態の表示を統一する
 */
export function ChartContainer({
  title,
  isLoading = false,
  errorMessage,
  emptyMessage = 'データがありません',
  isEmpty = false,
  height = chartDefaults.height,
  children,
}: ChartContainerProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base font-semibold">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div style={{ height }} className="flex items-center justify-center">
            <Loading />
          </div>
        ) : errorMessage ? (
          <div
            style={{ height }}
            className="flex items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 text-sm text-red-800"
          >
            {errorMessage}
          </div>
        ) : isEmpty ? (
          <div
            style={{ height }}
            className="flex items-center justify-center text-sm text-muted-foreground"
          >
            {emptyMessage}
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={height}>
            {children as React.ReactElement}
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  );
}
