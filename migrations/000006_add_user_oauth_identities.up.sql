  CREATE TABLE todoapp.user_oauth_identities (
      id            SERIAL       PRIMARY KEY,
      user_id       INT          NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
      provider      VARCHAR(32)  NOT NULL CHECK (provider IN ('google', 'apple')),
      provider_sub  VARCHAR(255) NOT NULL,
      email         VARCHAR(255) NOT NULL,
      created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),

      UNIQUE (provider, provider_sub)
  );

  CREATE INDEX idx_user_oauth_identities_user_id
      ON todoapp.user_oauth_identities (user_id);
