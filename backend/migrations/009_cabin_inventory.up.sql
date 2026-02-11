CREATE TABLE IF NOT EXISTS cabin_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    voyage_id UUID NOT NULL REFERENCES voyages(id) ON DELETE CASCADE,
    cabin_type_id UUID NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
    total_cabins INTEGER NOT NULL DEFAULT 0,
    available_cabins INTEGER NOT NULL DEFAULT 0,
    reserved_cabins INTEGER NOT NULL DEFAULT 0,
    booked_cabins INTEGER NOT NULL DEFAULT 0,
    locked_cabins INTEGER NOT NULL DEFAULT 0,
    maintenance_cabins INTEGER NOT NULL DEFAULT 0,
    lock_version INTEGER DEFAULT 0,
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(voyage_id, cabin_type_id)
);

CREATE INDEX idx_inventory_voyage_id ON cabin_inventory(voyage_id);
CREATE INDEX idx_inventory_cabin_type_id ON cabin_inventory(cabin_type_id);
CREATE INDEX idx_inventory_available ON cabin_inventory(voyage_id, cabin_type_id) WHERE available_cabins > 0;

COMMENT ON TABLE cabin_inventory IS '舱房库存表，用于乐观锁控制并发预订';
COMMENT ON COLUMN cabin_inventory.lock_version IS '乐观锁版本号，用于并发控制';
