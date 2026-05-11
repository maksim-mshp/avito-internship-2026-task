CREATE TABLE IF NOT EXISTS ai_assistants_catalog.assistant_favorites (
    user_id UUID NOT NULL REFERENCES ai_assistants_catalog.users(id) ON DELETE CASCADE,
    assistant_id UUID NOT NULL REFERENCES ai_assistants_catalog.assistants(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, assistant_id)
);

CREATE INDEX IF NOT EXISTS idx_assistant_favorites_user_created_at
    ON ai_assistants_catalog.assistant_favorites (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_assistant_favorites_assistant_id
    ON ai_assistants_catalog.assistant_favorites (assistant_id);
