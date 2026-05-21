ALTER TABLE todoapp.users ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE todoapp.email_confirmation_tokens (
    token_hash  TEXT        PRIMARY KEY,
    user_id     INT         NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);