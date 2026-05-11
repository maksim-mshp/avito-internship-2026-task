ALTER TABLE ai_assistants_catalog.assistant_runs
    ADD COLUMN IF NOT EXISTS rating TEXT;

ALTER TABLE ai_assistants_catalog.assistant_runs
    ADD CONSTRAINT chk_assistant_runs_rating
    CHECK (rating IS NULL OR rating IN ('like', 'dislike'));

CREATE INDEX IF NOT EXISTS idx_assistant_runs_user_rating
    ON ai_assistants_catalog.assistant_runs (user_id, rating)
    WHERE rating IS NOT NULL;
