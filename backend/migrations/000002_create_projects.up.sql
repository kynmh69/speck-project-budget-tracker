-- Create projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'planning',
    budget_amount DECIMAL(15,2),
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT projects_status_check CHECK (status IN ('planning', 'in_progress', 'completed', 'on_hold')),
    CONSTRAINT projects_date_check CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

-- Indexes
CREATE INDEX projects_user_id_idx ON projects(user_id);
CREATE INDEX projects_status_idx ON projects(status);
CREATE INDEX projects_dates_idx ON projects(start_date, end_date);
CREATE INDEX projects_deleted_at_idx ON projects(deleted_at);
CREATE INDEX projects_name_idx ON projects(name);

-- Comments
COMMENT ON TABLE projects IS 'プロジェクト';
COMMENT ON COLUMN projects.status IS 'ステータス: planning, in_progress, completed, on_hold';
