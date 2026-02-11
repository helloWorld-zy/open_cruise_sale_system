CREATE TABLE IF NOT EXISTS cabin_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    voyage_id UUID NOT NULL REFERENCES voyages(id) ON DELETE CASCADE,
    cabin_type_id UUID NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
    price_type VARCHAR(20) NOT NULL DEFAULT 'standard',
    adult_price DECIMAL(10,2) NOT NULL,
    child_price DECIMAL(10,2),
    infant_price DECIMAL(10,2),
    single_supplement DECIMAL(10,2),
    port_fee DECIMAL(10,2) DEFAULT 0,
    service_fee DECIMAL(10,2) DEFAULT 0,
    is_promotion BOOLEAN DEFAULT false,
    promotion_start_date DATE,
    promotion_end_date DATE,
    min_passengers INTEGER DEFAULT 1,
    max_passengers INTEGER DEFAULT 4,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(voyage_id, cabin_type_id, price_type)
);

CREATE INDEX idx_prices_voyage_id ON cabin_prices(voyage_id);
CREATE INDEX idx_prices_cabin_type_id ON cabin_prices(cabin_type_id);
CREATE INDEX idx_prices_type ON cabin_prices(price_type);
CREATE INDEX idx_prices_promotion ON cabin_prices(is_promotion, promotion_start_date, promotion_end_date);

COMMENT ON TABLE cabin_prices IS '舱房价格表，支持多种价格类型';
COMMENT ON COLUMN cabin_prices.price_type IS '价格类型: standard-标准价, early_bird-早鸟价, last_minute-最后一分钟, group-团队价';
