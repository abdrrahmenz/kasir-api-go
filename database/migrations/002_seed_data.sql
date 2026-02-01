-- Seed categories
INSERT INTO categories (id, name, description) VALUES
(1, 'Makanan', 'Berbagai jenis makanan instan dan siap saji'),
(2, 'Minuman', 'Aneka minuman segar dan kemasan'),
(3, 'Bumbu Dapur', 'Perlengkapan bumbu untuk memasak')
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  updated_at = CURRENT_TIMESTAMP;

-- Seed products
INSERT INTO products (id, name, price, stock, category_id) VALUES
(1, 'Indomie Godog', 3500, 10, 1),
(2, 'Vit 1000ml', 3000, 40, 2),
(3, 'kecap', 12000, 20, 3)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  price = EXCLUDED.price,
  stock = EXCLUDED.stock,
  category_id = EXCLUDED.category_id,
  updated_at = CURRENT_TIMESTAMP;

-- Reset sequence for categories
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));

-- Reset sequence for products
SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));
