-- +goose Up
ALTER TABLE orders
    ADD COLUMN status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending','awaiting_payment','paid','fulfilled','shipped','delivered','cancelled','refunded')),
    ADD COLUMN currency TEXT NOT NULL DEFAULT 'usd',
    ADD COLUMN subtotal_cents INTEGER NOT NULL DEFAULT 0 CHECK (subtotal_cents >= 0),
    ADD COLUMN shipping_cents INTEGER NOT NULL DEFAULT 0 CHECK (shipping_cents >= 0),
    ADD COLUMN tax_cents INTEGER NOT NULL DEFAULT 0 CHECK (tax_cents >= 0),
    ADD COLUMN total_cents INTEGER NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    ADD COLUMN stripe_payment_intent_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN shipping_address_id BIGINT REFERENCES addresses(id),
    -- Immutable shipping-address snapshot: a copy taken at order time so that
    -- later edits/deletes of the source address never mutate historical orders.
    ADD COLUMN ship_recipient_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_line1 TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_line2 TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_city TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_region TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_postal_code TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_country TEXT NOT NULL DEFAULT '',
    ADD COLUMN ship_phone TEXT NOT NULL DEFAULT '',
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Immutable line-item snapshot: keep order history readable even if the
-- catalog (variant sku / product name) changes or is deleted later.
ALTER TABLE order_items
    ADD COLUMN variant_sku TEXT NOT NULL DEFAULT '',
    ADD COLUMN product_name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE order_items
    DROP COLUMN variant_sku,
    DROP COLUMN product_name;

ALTER TABLE orders
    DROP COLUMN status,
    DROP COLUMN currency,
    DROP COLUMN subtotal_cents,
    DROP COLUMN shipping_cents,
    DROP COLUMN tax_cents,
    DROP COLUMN total_cents,
    DROP COLUMN stripe_payment_intent_id,
    DROP COLUMN shipping_address_id,
    DROP COLUMN ship_recipient_name,
    DROP COLUMN ship_line1,
    DROP COLUMN ship_line2,
    DROP COLUMN ship_city,
    DROP COLUMN ship_region,
    DROP COLUMN ship_postal_code,
    DROP COLUMN ship_country,
    DROP COLUMN ship_phone,
    DROP COLUMN updated_at;
