-- Migration: Create refund_requests table
-- Up Migration

CREATE TABLE refund_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    order_item_id UUID REFERENCES order_items(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Refund amounts
    refund_amount DECIMAL(10, 2) NOT NULL,
    refund_reason TEXT NOT NULL,
    refund_type VARCHAR(20) NOT NULL DEFAULT 'partial', -- full, partial
    
    -- Status workflow
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, approved, rejected, processing, completed, failed
    
    -- Approval workflow
    requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    review_note TEXT,
    
    -- Processing
    processed_at TIMESTAMP,
    payment_refund_id VARCHAR(100),
    third_party_refund_id VARCHAR(100),
    
    -- Bank/Account info for refund (if not original payment method)
    refund_method VARCHAR(20) DEFAULT 'original', -- original, bank, alipay, wechat
    bank_name VARCHAR(100),
    bank_account VARCHAR(100),
    account_holder VARCHAR(100),
    
    -- Cancellation reason
    cancellation_reason VARCHAR(50), -- customer_request, voyage_cancelled, cabin_upgrade, other
    
    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT refund_requests_status_check CHECK (status IN ('pending', 'approved', 'rejected', 'processing', 'completed', 'failed')),
    CONSTRAINT refund_requests_type_check CHECK (refund_type IN ('full', 'partial')),
    CONSTRAINT refund_requests_method_check CHECK (refund_method IN ('original', 'bank', 'alipay', 'wechat')),
    CONSTRAINT refund_requests_reason_check CHECK (cancellation_reason IN ('customer_request', 'voyage_cancelled', 'cabin_upgrade', 'other'))
);

-- Indexes
CREATE INDEX idx_refund_requests_order_id ON refund_requests(order_id);
CREATE INDEX idx_refund_requests_user_id ON refund_requests(user_id);
CREATE INDEX idx_refund_requests_status ON refund_requests(status);
CREATE INDEX idx_refund_requests_requested_at ON refund_requests(requested_at DESC);
CREATE INDEX idx_refund_requests_reviewed_by ON refund_requests(reviewed_by);

-- Comments
COMMENT ON TABLE refund_requests IS '退款申请记录表';
COMMENT ON COLUMN refund_requests.order_id IS '关联订单ID';
COMMENT ON COLUMN refund_requests.refund_amount IS '退款金额';
COMMENT ON COLUMN refund_requests.refund_reason IS '退款原因描述';
COMMENT ON COLUMN refund_requests.status IS '退款状态: pending-待审核, approved-已批准, rejected-已拒绝, processing-处理中, completed-已完成, failed-失败';
COMMENT ON COLUMN refund_requests.reviewed_by IS '审核人ID';
COMMENT ON COLUMN refund_requests.review_note IS '审核备注';
COMMENT ON COLUMN refund_requests.cancellation_reason IS '取消原因分类';
