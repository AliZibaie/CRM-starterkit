CREATE TABLE IF NOT EXISTS cron_jobs (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    condition_json TEXT NOT NULL,
    action_json TEXT NOT NULL,
    cron_expr VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS cron_job_logs (
    id VARCHAR(36) PRIMARY KEY,
    job_id VARCHAR(36) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL,
    affected_rows BIGINT NOT NULL DEFAULT 0,
    error_message TEXT
);
