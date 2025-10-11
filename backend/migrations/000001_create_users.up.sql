-- Create users table
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'manager', 'member'))
);

-- Indexes
CREATE UNIQUE INDEX users_email_unique ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX users_role_idx ON users(role);
CREATE INDEX users_deleted_at_idx ON users(deleted_at);

-- Comments
COMMENT ON TABLE users IS 'システムユーザー';
COMMENT ON COLUMN users.email IS 'メールアドレス（ログインID）';
COMMENT ON COLUMN users.password_hash IS 'bcryptハッシュ化されたパスワード';
COMMENT ON COLUMN users.role IS 'ロール: admin, manager, member';
