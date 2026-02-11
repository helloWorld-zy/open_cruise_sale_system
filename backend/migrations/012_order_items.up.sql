CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    cabin_id UUID NOT NULL REFERENCES cabins(id),
    cabin_type_id UUID NOT NULL REFERENCES cabin_types(id),
    voyage_id UUID NOT NULL REFERENCES voyages(id),
    cabin_number VARCHAR(20),
    price_snapshot DECIMAL(10,2) NOT NULL,
    adult_count INTEGER NOT NULL DEFAULT 2,
    child_count INTEGER DEFAULT 0,
    infant_count INTEGER DEFAULT 0,
    adult_price DECIMAL(10,2) NOT NULL,
    child_price DECIMAL(10,2),
    infant_price DECIMAL(10,2),
    port_fee DECIMAL(10,2) DEFAULT 0,
    service_fee DECIMAL(10,2) DEFAULT 0,
    subtotal DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'confirmed',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_cabin_id ON order_items(cabin_id);
CREATE INDEX idx_order_items_status ON order_items(status);

COMMENT ON TABLE order_items IS '订单明细表';
COMMENT ON COLUMN order_items.price_snapshot IS '预订时的价格快照';
COMMENT ON COLUMN order_items.status IS '明细状态: confirmed-已确认, cancelled-已取消, changed-已变更';
