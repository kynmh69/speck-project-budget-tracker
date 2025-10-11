-- Create tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    assigned_to UUID,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    planned_hours DECIMAL(10,2) DEFAULT 0.00,
    actual_hours DECIMAL(10,2) DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'todo',
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT tasks_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT tasks_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES members(id) ON DELETE SET NULL,
    CONSTRAINT tasks_status_check CHECK (status IN ('todo', 'in_progress', 'completed', 'blocked')),
    CONSTRAINT tasks_planned_hours_check CHECK (planned_hours >= 0),
    CONSTRAINT tasks_actual_hours_check CHECK (actual_hours >= 0)
);

-- Indexes
CREATE INDEX tasks_project_id_idx ON tasks(project_id);
CREATE INDEX tasks_assigned_to_idx ON tasks(assigned_to);
CREATE INDEX tasks_status_idx ON tasks(status);
CREATE INDEX tasks_deleted_at_idx ON tasks(deleted_at);

-- Comments
COMMENT ON TABLE tasks IS 'タスク';
COMMENT ON COLUMN tasks.planned_hours IS '予定工数（時間）';
COMMENT ON COLUMN tasks.actual_hours IS '実績工数（時間）';
COMMENT ON COLUMN tasks.status IS 'ステータス: todo, in_progress, completed, blocked';
