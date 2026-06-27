-- name: ListProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductBySlug :one
SELECT * FROM products WHERE slug = $1;

-- name: ListActiveProducts :many
SELECT * FROM products WHERE status = 'active';

-- name: CreateProduct :one
INSERT INTO products (name, slug, description, status, category_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateCategory :one
INSERT INTO categories (name, slug, parent_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;

-- name: GetCategoryBySlug :one
SELECT * FROM categories WHERE slug = $1;

-- name: CreateProductOption :one
INSERT INTO product_options (product_id, name, position)
VALUES ($1, $2, $3) RETURNING *;

-- name: CreateProductOptionValue :one
INSERT INTO product_option_values (option_id, value, position)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListOptionsByProduct :many
SELECT * FROM product_options WHERE product_id = $1 ORDER BY position;

-- name: CreateProductVariant :one
INSERT INTO product_variants (product_id, sku, price_in_cents, stock, weight_grams)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetVariantByID :one
SELECT * FROM product_variants WHERE id = $1;

-- name: ListVariantsByProduct :many
SELECT * FROM product_variants WHERE product_id = $1 ORDER BY id;

-- name: DecreaseVariantStock :exec
UPDATE product_variants
SET stock = stock - $1
WHERE id = $2 AND stock >= $1;

-- name: LinkVariantOptionValue :exec
INSERT INTO variant_option_values (variant_id, option_value_id)
VALUES ($1, $2);

-- name: CreateProductImage :one
INSERT INTO product_images (product_id, variant_id, url, position)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListImagesByProduct :many
SELECT * FROM product_images WHERE product_id = $1 ORDER BY position;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, variant_id, quantity, price_in_cents)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmailIgnoreCase :one
SELECT * FROM users WHERE lower(email) = lower($1);

-- name: GetUserByUsernameIgnoreCase :one
SELECT * FROM users WHERE lower(username) = lower($1);

-- name: CreateUser :one
INSERT INTO users (id, username, email, hashed_password)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateAddress :one
INSERT INTO addresses (
    user_id, recipient_name, line1, line2, city, region, postal_code, country, phone, is_default
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: ListAddressesByUser :many
SELECT * FROM addresses
WHERE user_id = $1
ORDER BY is_default DESC, id;

-- name: GetAddressByIDForUser :one
SELECT * FROM addresses
WHERE id = $1 AND user_id = $2;

-- name: UpdateAddressForUser :one
UPDATE addresses
SET recipient_name = $3,
    line1 = $4,
    line2 = $5,
    city = $6,
    region = $7,
    postal_code = $8,
    country = $9,
    phone = $10,
    is_default = $11,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteAddressForUser :execrows
DELETE FROM addresses
WHERE id = $1 AND user_id = $2;

-- name: UnsetDefaultAddressesForUser :exec
UPDATE addresses
SET is_default = false, updated_at = now()
WHERE user_id = $1;

-- name: SetDefaultAddressForUser :exec
UPDATE addresses
SET is_default = true, updated_at = now()
WHERE id = $1 AND user_id = $2;
