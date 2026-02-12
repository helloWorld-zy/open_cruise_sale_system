CREATE TABLE travelogues (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    voyage_id BIGINT,
    cruise_id BIGINT,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    cover_image VARCHAR(500),
    images JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    destination_tags JSONB DEFAULT '[]',
    travel_date DATE,
    duration_days INT,
    companions VARCHAR(50), -- 'solo', 'couple', 'family', 'friends', 'group'
    rating DECIMAL(2,1), -- 1.0 to 5.0
    view_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    comment_count INT NOT NULL DEFAULT 0,
    share_count INT NOT NULL DEFAULT 0,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- 'draft', 'pending_review', 'published', 'rejected', 'archived'
    featured_order INT,
    meta_title VARCHAR(200),
    meta_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_travelogues_user_id ON travelogues(user_id);
CREATE INDEX idx_travelogues_voyage_id ON travelogues(voyage_id);
CREATE INDEX idx_travelogues_cruise_id ON travelogues(cruise_id);
CREATE INDEX idx_travelogues_status ON travelogues(status);
CREATE INDEX idx_travelogues_is_featured ON travelogues(is_featured) WHERE is_featured = TRUE;
CREATE INDEX idx_travelogues_is_published ON travelogues(is_published) WHERE is_published = TRUE;
CREATE INDEX idx_travelogues_published_at ON travelogues(published_at DESC);
CREATE INDEX idx_travelogues_view_count ON travelogues(view_count DESC);
CREATE INDEX idx_travelogues_rating ON travelogues(rating DESC);

-- Full-text search index (if using PostgreSQL's tsvector)
-- CREATE INDEX idx_travelogues_search ON travelogues USING GIN (to_tsvector('chinese', title || ' ' || content));

-- Foreign keys
ALTER TABLE travelogues ADD CONSTRAINT fk_travelogues_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Comments table for travelogues
CREATE TABLE travelogue_comments (
    id BIGSERIAL PRIMARY KEY,
    travelogue_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    parent_id BIGINT, -- For nested comments
    content TEXT NOT NULL,
    like_count INT NOT NULL DEFAULT 0,
    is_approved BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_travelogue_comments_travelogue_id ON travelogue_comments(travelogue_id);
CREATE INDEX idx_travelogue_comments_user_id ON travelogue_comments(user_id);
CREATE INDEX idx_travelogue_comments_parent_id ON travelogue_comments(parent_id);
CREATE INDEX idx_travelogue_comments_created_at ON travelogue_comments(created_at DESC);

ALTER TABLE travelogue_comments ADD CONSTRAINT fk_comments_travelogue 
    FOREIGN KEY (travelogue_id) REFERENCES travelogues(id) ON DELETE CASCADE;
ALTER TABLE travelogue_comments ADD CONSTRAINT fk_comments_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE travelogue_comments ADD CONSTRAINT fk_comments_parent 
    FOREIGN KEY (parent_id) REFERENCES travelogue_comments(id) ON DELETE CASCADE;

-- Likes table for travelogues
CREATE TABLE travelogue_likes (
    id BIGSERIAL PRIMARY KEY,
    travelogue_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(travelogue_id, user_id)
);

CREATE INDEX idx_travelogue_likes_travelogue_id ON travelogue_likes(travelogue_id);
CREATE INDEX idx_travelogue_likes_user_id ON travelogue_likes(user_id);

ALTER TABLE travelogue_likes ADD CONSTRAINT fk_likes_travelogue 
    FOREIGN KEY (travelogue_id) REFERENCES travelogues(id) ON DELETE CASCADE;
ALTER TABLE travelogue_likes ADD CONSTRAINT fk_likes_user 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
