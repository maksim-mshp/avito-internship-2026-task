CREATE TABLE IF NOT EXISTS ai_assistants_catalog.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO ai_assistants_catalog.users (id, email, role)
VALUES
    ('ca2a62f3-c998-4050-96c1-0c0f62cf6568', 'admin@example.com', 'admin'),
    ('44c75af3-eca3-4867-85fc-b8245eaafa3a', 'user@example.com', 'user')
ON CONFLICT (id) DO UPDATE
SET
    email = EXCLUDED.email,
    role = EXCLUDED.role;

ALTER TABLE ai_assistants_catalog.assistant_runs
    ADD CONSTRAINT fk_assistant_runs_user_id
    FOREIGN KEY (user_id)
    REFERENCES ai_assistants_catalog.users(id);

CREATE INDEX IF NOT EXISTS idx_users_role ON ai_assistants_catalog.users (role);
