ALTER TABLE ai_assistants_catalog.users
    ADD COLUMN IF NOT EXISTS password_hash TEXT;

UPDATE ai_assistants_catalog.users
SET password_hash = '$2a$10$hfcLxclybGOmg8CIJeuZi.3qIHtMKLdw2bHne4ZRDbARnai0Sgh96'
WHERE password_hash IS NULL;

ALTER TABLE ai_assistants_catalog.users
    ALTER COLUMN password_hash SET NOT NULL;
