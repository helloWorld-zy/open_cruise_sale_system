CREATE TABLE IF NOT EXISTS cabins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    voyage_id UUID NOT NULL REFERENCES voyages(id) ON DELETE CASCADE,
    cabin_type_id UUID NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
    cabin_number VARCHAR(20) NOT NULL,
    deck_number INTEGER,
    section VARCHAR(10),
    status VARCHAR(20) DEFAULT 'available',
    is_accessible BOOLEAN DEFAULT false,
    is_connecting BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(voyage_id, cabin_number)
);

CREATE INDEX idx_cabins_voyage_id ON cabins(voyage_id);
CREATE INDEX idx_cabins_cabin_type_id ON cabins(cabin_type_id);
CREATE INDEX idx_cabins_status ON cabins(status);
CREATE INDEX idx_cabins_deleted_at ON cabins(deleted_at) WHERE deleted_at IS NULL;

COMMENT ON TABLE cabins IS '具体舱房实例（每航次每个舱房一条记录）';
COMMENT ON COLUMN cabins.status IS '舱房状态: available-可预订, occupied-已占用, maintenance-维护中, locked-已锁定';
