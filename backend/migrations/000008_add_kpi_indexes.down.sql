-- KPI集計クエリ最適化用インデックスの削除 (T164)
DROP INDEX IF EXISTS projects_user_status_idx;
DROP INDEX IF EXISTS projects_user_updated_idx;
DROP INDEX IF EXISTS budgets_profit_rate_idx;
