CREATE TABLE IF NOT EXISTS voyages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    cruise_id UUID NOT NULL REFERENCES cruises(id) ON DELETE CASCADE,
    voyage_number VARCHAR(50) NOT NULL UNIQUE,
    departure_date DATE NOT NULL,
    arrival_date DATE NOT NULL,
    departure_time TIME,
    arrival_time TIME,
    status VARCHAR(20) DEFAULT 'scheduled',
    booking_status VARCHAR(20) DEFAULT 'open',
    min_price DECIMAL(10,2),
    max_price DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_voyages_route_id ON voyages(route_id);
CREATE INDEX idx_voyages_cruise_id ON voyages(cruise_id);
CREATE INDEX idx_voyages_departure_date ON voyages(departure_date);
CREATE INDEX idx_voyages_status ON voyages(status);
CREATE INDEX idx_voyages_booking_status ON voyages(booking_status);
CREATE INDEX idx_voyages_deleted_at ON voyages(deleted_at) WHERE deleted_at IS NULL;

COMMENT ON TABLE voyages IS '邮轮航次信息';
COMMENT ON COLUMN voyages.status IS '航次状态: scheduled-计划中, active-进行中, completed-已完成, cancelled-已取消';
COMMENT ON COLUMN voyages.booking_status IS '预订状态: open-开放预订, full-已满, closed-已关闭';
