CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    payment_no VARCHAR(100) NOT NULL UNIQUE,
    payment_method VARCHAR(20) NOT NULL,
    payment_channel VARCHAR(50),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'CNY',
    status VARCHAR(20) DEFAULT 'pending',
    third_party_transaction_id VARCHAR(200),
    third_party_response TEXT,
    paid_at TIMESTAMP WITH TIME ZONE,
    notify_at TIMESTAMP WITH TIME ZONE,
    notify_data JSONB,
    retry_count INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_payment_no ON payments(payment_no);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_third_party_id ON payments(third_party_transaction_id);

COMMENT ON TABLE payments IS '支付记录表';
COMMENT ON COLUMN payments.payment_method IS '支付方式: wechat-微信支付, alipay-支付宝, card-银行卡';
COMMENT ON COLUMN payments.status IS '支付状态: pending-待支付, processing-处理中, success-成功, failed-失败, cancelled-已取消';
COMMENT ON COLUMN payments.third_party_transaction_id IS '第三方支付交易号';
COMMENT ON COLUMN payments.notify_data IS '支付回调通知数据';
