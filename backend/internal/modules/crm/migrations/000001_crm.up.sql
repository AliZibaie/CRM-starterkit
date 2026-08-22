CREATE TABLE IF NOT EXISTS crm_leads (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL,
    owner_id VARCHAR(36) NOT NULL DEFAULT '',
    source VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_contacts (
    id VARCHAR(36) PRIMARY KEY,
    lead_id VARCHAR(36) NOT NULL DEFAULT '',
    user_id VARCHAR(36) NOT NULL DEFAULT '',
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    owner_id VARCHAR(36) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_pipelines (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stages_json TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_deals (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    contact_id VARCHAR(36) NOT NULL DEFAULT '',
    pipeline_id VARCHAR(36) NOT NULL DEFAULT '',
    stage VARCHAR(100) NOT NULL DEFAULT '',
    owner_id VARCHAR(36) NOT NULL DEFAULT '',
    amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_activities (
    id VARCHAR(36) PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    body TEXT NOT NULL,
    actor_id VARCHAR(36) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    entity_type VARCHAR(50) NOT NULL DEFAULT '',
    entity_id VARCHAR(36) NOT NULL DEFAULT '',
    owner_id VARCHAR(36) NOT NULL,
    status VARCHAR(50) NOT NULL,
    due_at TIMESTAMP NULL,
    reminder_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_tags (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);
CREATE TABLE IF NOT EXISTS crm_entity_tags (
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    tag_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (entity_type, entity_id, tag_id)
);
CREATE TABLE IF NOT EXISTS crm_segments (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    rules_json TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
