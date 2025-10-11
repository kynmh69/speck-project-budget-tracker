-- Create time_entries table
CREATE TABLE time_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    member_id UUID NOT NULL,
    user_id UUID NOT NULL,
    work_date DATE NOT NULL DEFAULT CURRENT_DATE,
    hours DECIMAL(5,2) NOT NULL,
    hourly_rate_snapshot DECIMAL(10,2),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT time_entries_task_id_fkey FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT time_entries_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE RESTRICT,
    CONSTRAINT time_entries_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT time_entries_hours_check CHECK (hours > 0 AND hours <= 24)
);

-- Indexes
CREATE INDEX time_entries_task_id_idx ON time_entries(task_id);
CREATE INDEX time_entries_member_id_idx ON time_entries(member_id);
CREATE INDEX time_entries_user_id_idx ON time_entries(user_id);
CREATE INDEX time_entries_work_date_idx ON time_entries(work_date);
CREATE INDEX time_entries_task_date_idx ON time_entries(task_id, work_date);

-- Comments
COMMENT ON TABLE time_entries IS '工数記録';
COMMENT ON COLUMN time_entries.hours IS '工数（時間）';
COMMENT ON COLUMN time_entries.hourly_rate_snapshot IS '記録時点の単価';
