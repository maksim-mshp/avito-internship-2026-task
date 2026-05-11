DROP INDEX IF EXISTS ai_assistants_catalog.idx_assistant_runs_user_rating;

ALTER TABLE ai_assistants_catalog.assistant_runs
    DROP CONSTRAINT IF EXISTS chk_assistant_runs_rating;

ALTER TABLE ai_assistants_catalog.assistant_runs
    DROP COLUMN IF EXISTS rating;
