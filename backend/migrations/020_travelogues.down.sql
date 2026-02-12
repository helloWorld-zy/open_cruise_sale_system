ALTER TABLE travelogue_likes DROP CONSTRAINT IF EXISTS fk_likes_travelogue;
ALTER TABLE travelogue_likes DROP CONSTRAINT IF EXISTS fk_likes_user;
DROP INDEX IF EXISTS idx_travelogue_likes_user_id;
DROP INDEX IF EXISTS idx_travelogue_likes_travelogue_id;
DROP TABLE IF EXISTS travelogue_likes;

ALTER TABLE travelogue_comments DROP CONSTRAINT IF EXISTS fk_comments_parent;
ALTER TABLE travelogue_comments DROP CONSTRAINT IF EXISTS fk_comments_user;
ALTER TABLE travelogue_comments DROP CONSTRAINT IF EXISTS fk_comments_travelogue;
DROP INDEX IF EXISTS idx_travelogue_comments_created_at;
DROP INDEX IF EXISTS idx_travelogue_comments_parent_id;
DROP INDEX IF EXISTS idx_travelogue_comments_user_id;
DROP INDEX IF EXISTS idx_travelogue_comments_travelogue_id;
DROP TABLE IF EXISTS travelogue_comments;

ALTER TABLE travelogues DROP CONSTRAINT IF EXISTS fk_travelogues_user;
DROP INDEX IF EXISTS idx_travelogues_rating;
DROP INDEX IF EXISTS idx_travelogues_view_count;
DROP INDEX IF EXISTS idx_travelogues_published_at;
DROP INDEX IF EXISTS idx_travelogues_is_published;
DROP INDEX IF EXISTS idx_travelogues_is_featured;
DROP INDEX IF EXISTS idx_travelogues_status;
DROP INDEX IF EXISTS idx_travelogues_cruise_id;
DROP INDEX IF EXISTS idx_travelogues_voyage_id;
DROP INDEX IF EXISTS idx_travelogues_user_id;
DROP TABLE IF EXISTS travelogues;
