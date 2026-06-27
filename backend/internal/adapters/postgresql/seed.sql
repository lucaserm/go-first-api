-- Demo seed data for local development. Idempotent-ish: safe on a fresh DB.
INSERT INTO categories (name, slug) VALUES
    ('Apparel', 'apparel'),
    ('Accessories', 'accessories');

INSERT INTO products (name, slug, description, status, category_id) VALUES
    ('Field Jacket', 'field-jacket', 'A waxed-cotton field jacket built for weather. Four-pocket front, corduroy collar.', 'active', (SELECT id FROM categories WHERE slug = 'apparel')),
    ('Merino Tee', 'merino-tee', 'Lightweight 150gsm merino crew. Breathes in summer, layers in winter.', 'active', (SELECT id FROM categories WHERE slug = 'apparel')),
    ('Canvas Tote', 'canvas-tote', 'Heavy 18oz cotton canvas tote with riveted handles.', 'active', (SELECT id FROM categories WHERE slug = 'accessories'));

INSERT INTO product_variants (product_id, sku, price_in_cents, stock, weight_grams) VALUES
    ((SELECT id FROM products WHERE slug = 'field-jacket'), 'FJ-OLV-M', 18900, 12, 1200),
    ((SELECT id FROM products WHERE slug = 'field-jacket'), 'FJ-OLV-L', 18900,  5, 1250),
    ((SELECT id FROM products WHERE slug = 'field-jacket'), 'FJ-BLK-M', 18900,  0, 1200),
    ((SELECT id FROM products WHERE slug = 'merino-tee'),   'MT-GRY-S',  5900, 30,  180),
    ((SELECT id FROM products WHERE slug = 'merino-tee'),   'MT-GRY-M',  5900, 22,  190),
    ((SELECT id FROM products WHERE slug = 'merino-tee'),   'MT-NVY-M',  5900, 14,  190),
    ((SELECT id FROM products WHERE slug = 'canvas-tote'),  'CT-NAT-OS', 3400, 50,  400);

INSERT INTO product_options (product_id, name, position) VALUES
    ((SELECT id FROM products WHERE slug = 'field-jacket'), 'Size', 0),
    ((SELECT id FROM products WHERE slug = 'field-jacket'), 'Color', 1),
    ((SELECT id FROM products WHERE slug = 'merino-tee'), 'Size', 0),
    ((SELECT id FROM products WHERE slug = 'merino-tee'), 'Color', 1);
