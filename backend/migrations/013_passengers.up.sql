CREATE TABLE IF NOT EXISTS passengers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    order_item_id UUID NOT NULL REFERENCES order_items(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    given_name VARCHAR(100),
    gender VARCHAR(10) NOT NULL,
    birth_date DATE NOT NULL,
    nationality VARCHAR(50),
    passport_number VARCHAR(50),
    passport_expiry DATE,
    id_number VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    passenger_type VARCHAR(20) NOT NULL DEFAULT 'adult',
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(20),
    dietary_requirements TEXT,
    medical_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_passengers_order_id ON passengers(order_id);
CREATE INDEX idx_passengers_order_item_id ON passengers(order_item_id);
CREATE INDEX idx_passengers_passport ON passengers(passport_number);
CREATE INDEX idx_passengers_type ON passengers(passenger_type);

COMMENT ON TABLE passengers IS '乘客信息表';
COMMENT ON COLUMN passengers.passenger_type IS '乘客类型: adult-成人, child-儿童, infant-婴儿';
COMMENT ON COLUMN passengers.nationality IS '国籍';
