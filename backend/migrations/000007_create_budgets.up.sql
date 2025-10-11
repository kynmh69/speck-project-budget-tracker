-- Create budgets table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL UNIQUE,
    revenue DECIMAL(15,2) DEFAULT 0.00,
    total_cost DECIMAL(15,2) DEFAULT 0.00,
    profit DECIMAL(15,2) DEFAULT 0.00,
    profit_rate DECIMAL(5,2) DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'JPY',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT budgets_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT budgets_revenue_check CHECK (revenue >= 0),
    CONSTRAINT budgets_total_cost_check CHECK (total_cost >= 0)
);

-- Indexes
CREATE UNIQUE INDEX budgets_project_id_unique ON budgets(project_id);

-- Comments
COMMENT ON TABLE budgets IS 'プロジェクト予算・収支';
COMMENT ON COLUMN budgets.revenue IS '売上金額';
COMMENT ON COLUMN budgets.total_cost IS '総コスト（計算値）';
COMMENT ON COLUMN budgets.profit IS '利益（計算値）';
COMMENT ON COLUMN budgets.profit_rate IS '利益率（％、計算値）';
COMMENT ON COLUMN budgets.currency IS '通貨（ISO 4217）';
