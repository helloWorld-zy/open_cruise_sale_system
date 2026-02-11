CREATE TABLE IF NOT EXISTS facility_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cruise_id UUID NOT NULL REFERENCES cruises(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(255),
    sort_weight INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_facility_categories_cruise_id ON facility_categories(cruise_id, sort_weight) WHERE deleted_at IS NULL;

COMMENT ON TABLE facility_categories IS '设施分类表';
COMMENT ON COLUMN facility_categories.cruise_id IS '所属邮轮ID';
COMMENT ON COLUMN facility_categories.name IS '分类名称';
COMMENT ON COLUMN facility_categories.icon IS '图标';
