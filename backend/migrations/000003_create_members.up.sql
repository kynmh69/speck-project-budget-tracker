-- Create members table
CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    hourly_rate DECIMAL(10,2) DEFAULT 0.00,
    department VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT members_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT members_hourly_rate_check CHECK (hourly_rate >= 0)
);

-- Indexes
CREATE INDEX members_user_id_idx ON members(user_id);
CREATE INDEX members_email_idx ON members(email);
CREATE INDEX members_deleted_at_idx ON members(deleted_at);

-- Comments
COMMENT ON TABLE members IS 'プロジェクトメンバー';
COMMENT ON COLUMN members.hourly_rate IS '時間単価';
