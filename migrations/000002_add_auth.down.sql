DROP TABLE IF EXISTS todoapp.refresh_tokens;

DROP INDEX IF EXISTS todoapp.users_email_unique;

ALTER TABLE todoapp.users
    DROP COLUMN IF EXISTS password_hash;
