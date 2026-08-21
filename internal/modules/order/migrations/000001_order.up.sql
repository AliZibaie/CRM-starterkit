CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    status VARCHAR(30) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    subtotal BIGINT NOT NULL,
    discount_code VARCHAR(50) NOT NULL DEFAULT '',
    discount_amount BIGINT NOT NULL DEFAULT 0,
    total_amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    item_id VARCHAR(36) NOT NULL,
    variant_id VARCHAR(36) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    unit_price BIGINT NOT NULL,
    quantity BIGINT NOT NULL,
    line_total BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL UNIQUE,
    number VARCHAR(50) NOT NULL UNIQUE,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    status VARCHAR(30) NOT NULL
);
