CREATE TABLE notification_settings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    order_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    payment_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    inventory_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    system_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    refund_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    voyage_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    promotion_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    wechat_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sms_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    email_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    quiet_hours_start VARCHAR(5), -- "HH:MM" format
    quiet_hours_end VARCHAR(5),   -- "HH:MM" format
    quiet_hours_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index
CREATE UNIQUE INDEX idx_notification_settings_user_id ON notification_settings(user_id);

-- Foreign key
ALTER TABLE notification_settings ADD CONSTRAINT fk_notification_settings_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
