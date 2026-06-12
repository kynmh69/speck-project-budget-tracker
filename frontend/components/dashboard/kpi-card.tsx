'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import { LucideIcon } from 'lucide-react';

interface KpiCardProps {
  /** カードのタイトル（例: 総プロジェクト数） */
  title: string;
  /** 表示する値（フォーマット済み文字列または数値） */
  value: string | number;
  /** タイトル横に表示するアイコン */
  icon?: LucideIcon;
  /** 値の下に表示する補足テキスト */
  description?: string;
  /** 値に適用する追加クラス（例: 赤字時の警告色） */
  valueClassName?: string;
}

export function KpiCard({ title, value, icon: Icon, description, valueClassName }: KpiCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        {Icon && <Icon className="h-4 w-4 text-muted-foreground" />}
      </CardHeader>
      <CardContent>
        <div className={cn('text-2xl font-bold', valueClassName)}>{value}</div>
        {description && <p className="text-xs text-muted-foreground mt-1">{description}</p>}
      </CardContent>
    </Card>
  );
}
