CREATE TABLE IF NOT EXISTS system_settings (
    setting_key VARCHAR(191) PRIMARY KEY,
    value TEXT,
    updated_by VARCHAR(36) NOT NULL DEFAULT '',
    updated_at TIMESTAMP NOT NULL
);
