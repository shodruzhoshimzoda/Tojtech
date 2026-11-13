ALTER TABLE users ADD COLUMN IF NOT EXISTS username VARCHAR(20);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower ON users (LOWER(username));