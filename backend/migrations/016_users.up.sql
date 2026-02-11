-- Migration: Create users table
-- Up Migration

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic info
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255),
    
    -- WeChat info
    wechat_openid VARCHAR(100) UNIQUE,
    wechat_unionid VARCHAR(100),
    wechat_nickname VARCHAR(100),
    wechat_avatar_url TEXT,
    
    -- Profile
    nickname VARCHAR(100),
    avatar_url TEXT,
    real_name VARCHAR(100),
    gender VARCHAR(10), -- male, female, unknown
    birth_date DATE,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, banned
    
    -- Verification
    phone_verified BOOLEAN DEFAULT FALSE,
    email_verified BOOLEAN DEFAULT FALSE,
    identity_verified BOOLEAN DEFAULT FALSE,
    
    -- Identity info
    id_number VARCHAR(18),
    id_front_image TEXT,
    id_back_image TEXT,
    
    -- Last login
    last_login_at TIMESTAMP,
    last_login_ip INET,
    
    -- Soft delete
    deleted_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive', 'banned')),
    CONSTRAINT users_gender_check CHECK (gender IN ('male', 'female', 'unknown'))
);

-- Indexes
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_wechat_openid ON users(wechat_openid);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Comments
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.wechat_openid IS '微信OpenID';
COMMENT ON COLUMN users.phone IS '手机号码';
COMMENT ON COLUMN users.status IS '用户状态: active-正常, inactive-未激活, banned-禁用';
