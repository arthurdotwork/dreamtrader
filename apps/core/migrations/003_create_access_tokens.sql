-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS auth_access_tokens
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    internal_id SERIAL UNIQUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() at time zone 'utc'),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() at time zone 'utc'),
    deleted_at  TIMESTAMP WITH TIME ZONE NULL,

    user_id UUID NOT NULL REFERENCES users(id),
    access_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_auth_access_tokens_access_token ON auth_access_tokens(access_token);

CREATE TABLE IF NOT EXISTS auth_refresh_tokens
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    internal_id SERIAL UNIQUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() at time zone 'utc'),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() at time zone 'utc'),
    deleted_at  TIMESTAMP WITH TIME ZONE NULL,

    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_auth_refresh_tokens_refresh_token ON auth_refresh_tokens(refresh_token);

---- create above / drop below ----

DROP TABLE IF EXISTS auth_access_tokens;
DROP TABLE IF EXISTS auth_refresh_tokens;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
