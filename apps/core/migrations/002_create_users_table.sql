-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    internal_id SERIAL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() at time zone 'utc'),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() at time zone 'utc'),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,

    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
