'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { updateRevenueSchema, UpdateRevenueFormData } from '@/schemas/budget-schema';
import { useUpdateRevenue } from '@/hooks/use-budget';
import { formatCurrency } from '@/types/budget';
import { Edit2 } from 'lucide-react';

interface RevenueFormProps {
  projectId: string;
  currentRevenue: number;
  currency: string;
}

export function RevenueForm({ projectId, currentRevenue, currency }: RevenueFormProps) {
  const [open, setOpen] = useState(false);
  const updateRevenue = useUpdateRevenue();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<UpdateRevenueFormData>({
    resolver: zodResolver(updateRevenueSchema),
    defaultValues: {
      revenue: currentRevenue,
      currency: currency,
    },
  });

  const onSubmit = async (data: UpdateRevenueFormData) => {
    try {
      await updateRevenue.mutateAsync({
        projectId,
        data,
      });
      setOpen(false);
      reset({ revenue: data.revenue, currency: data.currency });
    } catch {
      // Error is handled by the mutation
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>売上設定</CardTitle>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm">
              <Edit2 className="h-4 w-4 mr-2" />
              編集
            </Button>
          </DialogTrigger>
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{formatCurrency(currentRevenue, currency)}</div>
          <p className="text-sm text-muted-foreground mt-1">プロジェクトの売上金額</p>
        </CardContent>
      </Card>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>売上を編集</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="revenue">売上金額</Label>
            <Input
              id="revenue"
              type="number"
              step="0.01"
              min="0"
              {...register('revenue', { valueAsNumber: true })}
            />
            {errors.revenue && (
              <p className="text-sm text-red-500">{errors.revenue.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="currency">通貨</Label>
            <Input
              id="currency"
              maxLength={3}
              placeholder="JPY"
              {...register('currency')}
            />
            {errors.currency && (
              <p className="text-sm text-red-500">{errors.currency.message}</p>
            )}
          </div>

          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? '保存中...' : '保存'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
