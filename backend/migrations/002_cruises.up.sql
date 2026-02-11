CREATE TABLE IF NOT EXISTS cruises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES cruise_companies(id) ON DELETE RESTRICT,
    name_cn VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    code VARCHAR(50) NOT NULL UNIQUE,
    gross_tonnage INTEGER,
    passenger_capacity INTEGER,
    crew_count INTEGER,
    built_year INTEGER,
    renovated_year INTEGER,
    length_meters DECIMAL(8,2),
    width_meters DECIMAL(8,2),
    deck_count INTEGER,
    cover_images JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance')),
    sort_weight INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_cruises_company_id ON cruises(company_id);
CREATE INDEX idx_cruises_status_sort ON cruises(status, sort_weight DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_cruises_deleted_at ON cruises(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_cruises_code ON cruises(code) WHERE deleted_at IS NULL;

COMMENT ON TABLE cruises IS '邮轮信息表';
COMMENT ON COLUMN cruises.company_id IS '所属公司ID';
COMMENT ON COLUMN cruises.name_cn IS '中文名称';
COMMENT ON COLUMN cruises.code IS '唯一代码';
COMMENT ON COLUMN cruises.gross_tonnage IS '总吨位';
COMMENT ON COLUMN cruises.passenger_capacity IS '乘客容量';
COMMENT ON COLUMN cruises.status IS '状态: active-上架, inactive-下架, maintenance-维护中';
