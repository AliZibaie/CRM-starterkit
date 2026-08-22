CREATE TABLE IF NOT EXISTS announcements (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    audience VARCHAR(50) NOT NULL,
    audience_ref VARCHAR(100) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL,
    created_by VARCHAR(36) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    broadcast_at TIMESTAMP NULL
);
