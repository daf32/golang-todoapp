DROP TABLE IF EXISTS todoapp.email_confirmation_tokens;

ALTER TABLE todoapp.users DROP COLUMN IF EXISTS email_verified;