-- +goose Up
ALTER TABLE products ADD COLUMN slug TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS products_slug_key ON products (slug);
ALTER TABLE products ADD COLUMN description TEXT NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived'));
ALTER TABLE products ADD COLUMN category_id BIGINT REFERENCES categories(id);
ALTER TABLE products DROP COLUMN price_in_cents;
ALTER TABLE products DROP COLUMN quantity;

-- +goose Down
ALTER TABLE products ADD COLUMN quantity INTEGER NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN price_in_cents INTEGER NOT NULL DEFAULT 0 CHECK (price_in_cents >= 0);
ALTER TABLE products DROP COLUMN category_id;
ALTER TABLE products DROP COLUMN status;
ALTER TABLE products DROP COLUMN description;
DROP INDEX IF EXISTS products_slug_key;
ALTER TABLE products DROP COLUMN slug;
