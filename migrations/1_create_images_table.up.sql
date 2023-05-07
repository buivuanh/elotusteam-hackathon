CREATE TABLE IF NOT EXISTS public.images(
    id              SERIAL PRIMARY KEY,
    file_path       text NOT NULL UNIQUE,
    original_name   text NOT NULL,
    content_type    text NOT NULL,
    byte_size       integer NOT NULL,
    owner_id        integer NOT NULL,
    created_at      timestamp without time zone DEFAULT NOW(),
    deleted_at      timestamp without time zone,
    CONSTRAINT fk_owner_id FOREIGN KEY(owner_id) REFERENCES users(id)
);
