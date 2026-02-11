CREATE TABLE IF NOT EXISTS cruise_companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    logo_url VARCHAR(500),
    website VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_cruise_companies_name ON cruise_companies(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_cruise_companies_deleted_at ON cruise_companies(deleted_at) WHERE deleted_at IS NOT NULL;

COMMENT ON TABLE cruise_companies IS '邮轮公司信息表';
COMMENT ON COLUMN cruise_companies.id IS '公司ID';
COMMENT ON COLUMN cruise_companies.name IS '公司名称';
COMMENT ON COLUMN cruise_companies.name_en IS '英文名称';
COMMENT ON COLUMN cruise_companies.logo_url IS 'Logo URL';
COMMENT ON COLUMN cruise_companies.website IS '官网地址';
COMMENT ON COLUMN cruise_companies.description IS '公司简介';
