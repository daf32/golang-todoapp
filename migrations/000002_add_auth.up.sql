ALTER TABLE todoapp.users
    ADD COLUMN password_hash TEXT NOT NULL DEFAULT '';

ALTER TABLE todoapp.users
    ALTER COLUMN password_hash DROP DEFAULT;

CREATE UNIQUE INDEX users_email_unique
    ON todoapp.users (email);

CREATE TABLE todoapp.refresh_tokens (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         INTEGER         NOT NULL    REFERENCES todoapp.users(id) ON DELETE CASCADE,
    token           VARCHAR(255)    UNIQUE      NOT NULL,
    expires_at      TIMESTAMPTZ     NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    revoked         BOOLEAN         NOT NULL    DEFAULT FALSE
);

CREATE INDEX idx_refresh_tokens_token
    ON todoapp.refresh_tokens (token);
