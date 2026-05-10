CREATE SCHEMA IF NOT EXISTS ai_assistants_catalog;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS ai_assistants_catalog.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_created_at ON ai_assistants_catalog.categories (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_categories_name_lower ON ai_assistants_catalog.categories (LOWER(name));
