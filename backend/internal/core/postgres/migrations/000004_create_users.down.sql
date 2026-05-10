ALTER TABLE ai_assistants_catalog.assistant_runs
    DROP CONSTRAINT IF EXISTS fk_assistant_runs_user_id;

DROP TABLE IF EXISTS ai_assistants_catalog.users;
