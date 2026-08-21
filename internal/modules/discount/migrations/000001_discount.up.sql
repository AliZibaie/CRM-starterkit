CREATE TABLE IF NOT EXISTS discount_codes (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    value BIGINT NOT NULL,
    min_order_amount BIGINT NOT NULL DEFAULT 0,
    max_uses BIGINT NOT NULL DEFAULT 0,
    used_count BIGINT NOT NULL DEFAULT 0,
    starts_at TIMESTAMP NULL,
    ends_at TIMESTAMP NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
