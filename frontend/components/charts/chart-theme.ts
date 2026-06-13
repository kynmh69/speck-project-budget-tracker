// グラフ共通のカラーテーマとフォーマッタ (T179)
// 全チャートコンポーネントはここから色・書式を参照して統一感を保つ

export const chartColors = {
  planned: '#94a3b8', // slate-400: 予定
  actual: '#3b82f6', // blue-500: 実績
  revenue: '#3b82f6', // blue-500: 売上
  cost: '#f59e0b', // amber-500: コスト
  profit: '#22c55e', // green-500: 利益
  deficit: '#ef4444', // red-500: 赤字
  // 円グラフなど系列数が可変のグラフ用パレット
  series: [
    '#3b82f6', // blue-500
    '#22c55e', // green-500
    '#f59e0b', // amber-500
    '#a855f7', // purple-500
    '#ec4899', // pink-500
    '#14b8a6', // teal-500
    '#f97316', // orange-500
    '#6366f1', // indigo-500
  ],
} as const;

export const chartDefaults = {
  height: 320,
  margin: { top: 8, right: 16, bottom: 8, left: 16 },
  gridStroke: '#e2e8f0', // slate-200
  axisFontSize: 12,
} as const;

export function seriesColor(index: number): string {
  return chartColors.series[index % chartColors.series.length];
}

export function formatHours(value: number): string {
  return `${value.toLocaleString()}h`;
}

export function formatYen(value: number): string {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(value);
}

export function formatPercent(value: number): string {
  return `${value.toFixed(1)}%`;
}
