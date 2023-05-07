CREATE TABLE IF NOT EXISTS public.users(
    id              SERIAL PRIMARY KEY,
    username        text NOT NULL UNIQUE,
    hashed_password text NOT NULL,
    created_at      timestamp without time zone DEFAULT NOW(),
    deleted_at      timestamp without time zone
);
