-- +goose Up
ALTER TABLE users
    ADD COLUMN role TEXT NOT NULL DEFAULT 'customer'
        CHECK (role IN ('customer', 'admin'));

-- +goose Down
ALTER TABLE users DROP COLUMN role;
