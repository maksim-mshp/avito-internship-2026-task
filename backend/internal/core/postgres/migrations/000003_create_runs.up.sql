CREATE TABLE IF NOT EXISTS ai_assistants_catalog.assistant_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assistant_id UUID NOT NULL REFERENCES ai_assistants_catalog.assistants(id),
    user_id UUID NOT NULL,
    model VARCHAR(100) NOT NULL,
    user_prompt TEXT NOT NULL,
    output TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'success', 'failed')),
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_assistant_runs_user_created_at ON ai_assistants_catalog.assistant_runs (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_assistant_runs_assistant_created_at ON ai_assistants_catalog.assistant_runs (assistant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_assistant_runs_status ON ai_assistants_catalog.assistant_runs (status);
