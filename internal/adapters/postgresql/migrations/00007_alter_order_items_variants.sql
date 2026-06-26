-- +goose Up
ALTER TABLE order_items ADD COLUMN variant_id BIGINT NOT NULL REFERENCES product_variants(id);
ALTER TABLE order_items DROP COLUMN product_id;

-- +goose Down
ALTER TABLE order_items ADD COLUMN product_id BIGINT NOT NULL;
ALTER TABLE order_items DROP COLUMN variant_id;
