-- migrations/000006_create_refresh_tokens.up.sql
CREATE TABLE refresh_tokens (
    id          BIGSERIAL PRIMARY KEY,
    user_uuid   UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    token_hash  VARCHAR(64) NOT NULL UNIQUE, -- sha256 in hex = 64 symbols
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_uuid ON refresh_tokens(user_uuid);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);