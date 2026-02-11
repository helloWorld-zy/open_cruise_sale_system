CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    user_id UUID,
    voyage_id UUID NOT NULL REFERENCES voyages(id),
    cruise_id UUID NOT NULL REFERENCES cruises(id),
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    paid_amount DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'CNY',
    status VARCHAR(20) DEFAULT 'pending',
    payment_status VARCHAR(20) DEFAULT 'unpaid',
    passenger_count INTEGER NOT NULL DEFAULT 1,
    cabin_count INTEGER NOT NULL DEFAULT 1,
    contact_name VARCHAR(100),
    contact_phone VARCHAR(20),
    contact_email VARCHAR(100),
    remark TEXT,
    booked_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_voyage_id ON orders(voyage_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_expires_at ON orders(expires_at);
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at) WHERE deleted_at IS NULL;

COMMENT ON TABLE orders IS '订单主表';
COMMENT ON COLUMN orders.status IS '订单状态: pending-待支付, paid-已支付, confirmed-已确认, cancelled-已取消, completed-已完成, refunded-已退款';
COMMENT ON COLUMN orders.payment_status IS '支付状态: unpaid-未支付, partial-部分支付, paid-已支付, refunded-已退款';
COMMENT ON COLUMN orders.expires_at IS '订单超时时间，超过此时间未支付自动取消';
