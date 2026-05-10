CREATE TABLE IF NOT EXISTS ai_assistants_catalog.assistants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES ai_assistants_catalog.categories(id),
    name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    model VARCHAR(100) NOT NULL,
    system_prompt TEXT NOT NULL,
    example_user_prompt TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_assistants_category_id ON ai_assistants_catalog.assistants (category_id);
CREATE INDEX IF NOT EXISTS idx_assistants_active_created_at ON ai_assistants_catalog.assistants (is_active, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_assistants_name ON ai_assistants_catalog.assistants (name);
CREATE INDEX IF NOT EXISTS idx_assistants_description ON ai_assistants_catalog.assistants (description);
