CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'order', 'payment', 'inventory', 'system', etc.
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    data JSONB, -- Additional notification data (order_id, voyage_id, etc.)
    channel VARCHAR(20) NOT NULL DEFAULT 'in_app', -- 'in_app', 'wechat', 'sms', 'email'
    priority VARCHAR(10) NOT NULL DEFAULT 'normal', -- 'low', 'normal', 'high', 'urgent'
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    action_url VARCHAR(500), -- URL to navigate when clicking notification
    action_type VARCHAR(50), -- 'view_order', 'view_voyage', 'external_link', etc.
    source_id BIGINT, -- Reference to the source (order_id, voyage_id, etc.)
    source_type VARCHAR(50), -- 'order', 'voyage', 'system', etc.
    retry_count INT NOT NULL DEFAULT 0, -- For external channel delivery retries
    sent_at TIMESTAMP WITH TIME ZONE, -- When notification was sent (for external channels)
    delivered_at TIMESTAMP WITH TIME ZONE, -- When notification was delivered
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_user_read ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_user_archived ON notifications(user_id, is_archived);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_channel ON notifications(channel);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_priority ON notifications(priority);
CREATE INDEX idx_notifications_source ON notifications(source_type, source_id);

-- Composite index for unread notifications by user
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read, is_archived) WHERE is_read = FALSE AND is_archived = FALSE;

-- Foreign key
ALTER TABLE notifications ADD CONSTRAINT fk_notifications_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
