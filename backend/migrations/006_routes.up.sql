CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cruise_id UUID NOT NULL REFERENCES cruises(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    departure_port VARCHAR(100) NOT NULL,
    arrival_port VARCHAR(100) NOT NULL,
    duration_days INTEGER NOT NULL DEFAULT 1,
    description TEXT,
    itinerary JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'active',
    sort_weight INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_routes_cruise_id ON routes(cruise_id);
CREATE INDEX idx_routes_status ON routes(status);
CREATE INDEX idx_routes_code ON routes(code);
CREATE INDEX idx_routes_deleted_at ON routes(deleted_at) WHERE deleted_at IS NULL;

COMMENT ON TABLE routes IS '邮轮航线信息';
COMMENT ON COLUMN routes.itinerary IS '行程JSON数组，包含每日停靠港口和活动';
