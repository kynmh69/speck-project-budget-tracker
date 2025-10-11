-- Create project_members table
CREATE TABLE project_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    member_id UUID NOT NULL,
    joined_at DATE NOT NULL DEFAULT CURRENT_DATE,
    left_at DATE,
    allocation_rate DECIMAL(3,2) DEFAULT 1.00,
    hourly_rate_snapshot DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT project_members_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT project_members_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    CONSTRAINT project_members_allocation_rate_check CHECK (allocation_rate >= 0.0 AND allocation_rate <= 1.0),
    CONSTRAINT project_members_date_check CHECK (left_at IS NULL OR left_at >= joined_at)
);

-- Indexes
CREATE UNIQUE INDEX project_members_unique_idx ON project_members(project_id, member_id, joined_at);
CREATE INDEX project_members_project_id_idx ON project_members(project_id);
CREATE INDEX project_members_member_id_idx ON project_members(member_id);
CREATE INDEX project_members_dates_idx ON project_members(joined_at, left_at);

-- Comments
COMMENT ON TABLE project_members IS 'プロジェクトメンバー割り当て';
COMMENT ON COLUMN project_members.allocation_rate IS '割り当て率（0.0-1.0）';
COMMENT ON COLUMN project_members.hourly_rate_snapshot IS '参加時点の単価';
