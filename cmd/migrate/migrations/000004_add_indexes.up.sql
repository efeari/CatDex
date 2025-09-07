CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_cats_name ON cats USING gin (name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_cats_location ON cats USING gin (location gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_cats_created_at ON cats (created_at);

CREATE INDEX IF NOT EXISTS idx_users_name ON users (username);

CREATE INDEX IF NOT EXISTS idx_cats_user_id ON cats (user_id);