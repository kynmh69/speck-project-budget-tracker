-- KPI集計クエリ最適化用インデックス (T164)

-- ダッシュボードのステータス別件数集計（user_id + status の複合条件）
CREATE INDEX IF NOT EXISTS projects_user_status_idx
    ON projects(user_id, status)
    WHERE deleted_at IS NULL;

-- ダッシュボードの「最近のプロジェクト」（user_id絞り込み + updated_at降順）
CREATE INDEX IF NOT EXISTS projects_user_updated_idx
    ON projects(user_id, updated_at DESC)
    WHERE deleted_at IS NULL;

-- プロジェクト一覧の利益率フィルタ・ソート（budgets JOIN後の絞り込み）
CREATE INDEX IF NOT EXISTS budgets_profit_rate_idx
    ON budgets(profit_rate);
