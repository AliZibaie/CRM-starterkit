CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    actor_id VARCHAR(36) NOT NULL DEFAULT '',
    module VARCHAR(100) NOT NULL DEFAULT '',
    action VARCHAR(100) NOT NULL DEFAULT '',
    entity_type VARCHAR(100) NOT NULL DEFAULT '',
    entity_id VARCHAR(100) NOT NULL DEFAULT '',
    before_json TEXT,
    after_json TEXT,
    created_at TIMESTAMP NOT NULL
);
