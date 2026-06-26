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

-- name: GetUserByEmailIgnoreCase :one
SELECT * FROM users WHERE lower(email) = lower($1);

-- name: GetUserByUsernameIgnoreCase :one
SELECT * FROM users WHERE lower(username) = lower($1);

-- name: CreateUser :one
INSERT INTO users (id, username, email, hashed_password)
VALUES ($1, $2, $3, $4) RETURNING *;
