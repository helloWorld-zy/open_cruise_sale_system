ALTER TABLE notification_settings DROP CONSTRAINT IF EXISTS fk_notification_settings_user;
DROP INDEX IF EXISTS idx_notification_settings_user_id;
DROP TABLE IF EXISTS notification_settings;
