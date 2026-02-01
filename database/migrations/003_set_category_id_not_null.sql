-- Ensure category_id exists and enforce it as required on products
ALTER TABLE products
  ADD COLUMN IF NOT EXISTS category_id INTEGER;

-- Backfill existing rows (assumes category id 1 exists)
UPDATE products
  SET category_id = 1
  WHERE category_id IS NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'products_category_id_fkey'
  ) THEN
    ALTER TABLE products
      ADD CONSTRAINT products_category_id_fkey
      FOREIGN KEY (category_id) REFERENCES categories(id);
  END IF;
END $$;

ALTER TABLE products
  ALTER COLUMN category_id SET NOT NULL;
