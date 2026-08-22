CREATE TABLE IF NOT EXISTS cms_pages (
    id VARCHAR(36) PRIMARY KEY,
    slug VARCHAR(200) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS cms_blocks (
    id VARCHAR(36) PRIMARY KEY,
    page_id VARCHAR(36) NOT NULL,
    type VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS cms_menus (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL,
    items_json TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS cms_seo_meta (
    id VARCHAR(36) PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT,
    keywords VARCHAR(255) NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS cms_faqs (
    id VARCHAR(36) PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
