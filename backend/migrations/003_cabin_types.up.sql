CREATE TABLE IF NOT EXISTS cabin_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cruise_id UUID NOT NULL REFERENCES cruises(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    min_area_sqm DECIMAL(5,2),
    max_area_sqm DECIMAL(5,2),
    standard_guests INTEGER,
    max_guests INTEGER,
    bed_types VARCHAR(255),
    feature_tags JSONB DEFAULT '[]',
    description TEXT,
    images JSONB DEFAULT '[]',
    floor_plan_url VARCHAR(500),
    amenities JSONB DEFAULT '[]',
    sort_weight INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_cruise_cabin_code UNIQUE (cruise_id, code)
);

CREATE INDEX idx_cabin_types_cruise_id ON cabin_types(cruise_id);
CREATE INDEX idx_cabin_types_status_sort ON cabin_types(cruise_id, status, sort_weight) WHERE deleted_at IS NULL;
CREATE INDEX idx_cabin_types_deleted_at ON cabin_types(deleted_at) WHERE deleted_at IS NOT NULL;

COMMENT ON TABLE cabin_types IS '舱房类型表';
COMMENT ON COLUMN cabin_types.cruise_id IS '所属邮轮ID';
COMMENT ON COLUMN cabin_types.name IS '类型名称';
COMMENT ON COLUMN cabin_types.code IS '类型代码';
COMMENT ON COLUMN cabin_types.min_area_sqm IS '最小面积(平方米)';
COMMENT ON COLUMN cabin_types.max_area_sqm IS '最大面积(平方米)';
COMMENT ON COLUMN cabin_types.standard_guests IS '标准入住人数';
COMMENT ON COLUMN cabin_types.max_guests IS '最大入住人数';
