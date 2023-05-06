CREATE TABLE IF NOT EXISTS public.users (
    id         SERIAL PRIMARY KEY,
    username   text NOT NULL,
    password   text NOT NULL,
    created_at timestamp with time zone DEFAULT NOW(),
    deleted_at timestamp with time zone
);
