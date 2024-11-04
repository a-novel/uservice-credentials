CREATE TYPE credentials_role AS ENUM (
    'none',
    'early-access-program',
    'admin',
    'core'
);

--bun:split

CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    email TEXT NOT NULL,
    role credentials_role NOT NULL DEFAULT 'none',

    email_validation_token_id TEXT,
    pending_email_validation_token_id  TEXT,
    password_token_id  TEXT,
    reset_password_token_id TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    UNIQUE (email)
);
