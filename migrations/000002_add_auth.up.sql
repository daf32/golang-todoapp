-- Email/password auth from a clean slate: drop legacy users and dependent rows (e.g. tasks).
TRUNCATE TABLE todoapp.users RESTART IDENTITY CASCADE;

ALTER TABLE todoapp.users
    ADD COLUMN email VARCHAR(255) NOT NULL
        CHECK (
            char_length(email) BETWEEN 5 AND 255
            AND email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        ),
    ADD COLUMN password_hash TEXT NOT NULL;

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
