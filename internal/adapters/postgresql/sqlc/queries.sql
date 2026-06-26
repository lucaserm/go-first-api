-- name: ListProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_in_cents)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DecreaseProductStock :exec
UPDATE products
SET quantity = quantity - $1
WHERE id = $2 AND quantity >= $1;

-- name: GetUserByEmailIgnoreCase :one
SELECT * FROM users WHERE lower(email) = lower($1);

-- name: GetUserByUsernameIgnoreCase :one
SELECT * FROM users WHERE lower(username) = lower($1);

-- name: CreateUser :one
INSERT INTO users (id, username, email, hashed_password)
VALUES ($1, $2, $3, $4) RETURNING *;
