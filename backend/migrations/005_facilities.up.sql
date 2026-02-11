CREATE TABLE IF NOT EXISTS facilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cruise_id UUID NOT NULL REFERENCES cruises(id) ON DELETE CASCADE,
    category_id UUID REFERENCES facility_categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    deck_number INTEGER,
    open_time VARCHAR(255),
    is_free BOOLEAN DEFAULT true,
    price DECIMAL(10,2),
    description TEXT,
    images JSONB DEFAULT '[]',
    suitable_tags JSONB DEFAULT '[]',
    sort_weight INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'visible' CHECK (status IN ('visible', 'hidden')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_facilities_cruise_id ON facilities(cruise_id);
CREATE INDEX idx_facilities_category_id ON facilities(category_id);
CREATE INDEX idx_facilities_status ON facilities(cruise_id, category_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_facilities_deleted_at ON facilities(deleted_at) WHERE deleted_at IS NOT NULL;

COMMENT ON TABLE facilities IS '邮轮设施表';
COMMENT ON COLUMN facilities.cruise_id IS '所属邮轮ID';
COMMENT ON COLUMN facilities.category_id IS '分类ID';
COMMENT ON COLUMN facilities.name IS '设施名称';
COMMENT ON COLUMN facilities.deck_number IS '所在甲板层';
COMMENT ON COLUMN facilities.is_free IS '是否免费';
COMMENT ON COLUMN facilities.price IS '参考价格';
COMMENT ON COLUMN facilities.status IS '状态: visible-显示, hidden-隐藏';
