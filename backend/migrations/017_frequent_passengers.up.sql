-- Migration: Create frequent_passengers table
-- Up Migration

CREATE TABLE frequent_passengers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Passenger info
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    given_name VARCHAR(100),
    gender VARCHAR(10) NOT NULL, -- male, female
    birth_date DATE NOT NULL,
    nationality VARCHAR(50) DEFAULT '中国',
    
    -- Document info
    passport_number VARCHAR(50),
    passport_expiry DATE,
    id_number VARCHAR(18),
    
    -- Contact
    phone VARCHAR(20),
    email VARCHAR(100),
    
    -- Preferences
    dietary_requirements TEXT,
    medical_notes TEXT,
    
    -- Status
    is_default BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT frequent_passengers_gender_check CHECK (gender IN ('male', 'female')),
    CONSTRAINT frequent_passengers_document_check CHECK (
        (passport_number IS NOT NULL AND passport_number <> '') OR 
        (id_number IS NOT NULL AND id_number <> '')
    )
);

-- Indexes
CREATE INDEX idx_frequent_passengers_user_id ON frequent_passengers(user_id);
CREATE INDEX idx_frequent_passengers_passport ON frequent_passengers(passport_number);
CREATE INDEX idx_frequent_passengers_id_number ON frequent_passengers(id_number);

-- Comments
COMMENT ON TABLE frequent_passengers IS '常用乘客表';
COMMENT ON COLUMN frequent_passengers.is_default IS '是否为默认乘客';
