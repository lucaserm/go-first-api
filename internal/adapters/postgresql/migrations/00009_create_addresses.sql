-- +goose Up
CREATE TABLE IF NOT EXISTS addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_name TEXT NOT NULL,
    line1 TEXT NOT NULL,
    line2 TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL,
    region TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    country CHAR(2) NOT NULL CHECK (char_length(country) = 2),
    phone TEXT NOT NULL DEFAULT '',
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses (user_id);

-- Enforce at most one default address per user at the database level, so the
-- application-layer unset-then-set flow can't be defeated by a concurrent writer.
CREATE UNIQUE INDEX IF NOT EXISTS uq_addresses_one_default_per_user
    ON addresses (user_id) WHERE is_default;

-- +goose Down
DROP TABLE IF EXISTS addresses;
