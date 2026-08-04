DROP INDEX IF EXISTS idx_product_images_uuid;
DROP INDEX IF EXISTS idx_products_uuid;
DROP INDEX IF EXISTS idx_categories_uuid;

ALTER TABLE product_images DROP COLUMN IF EXISTS uuid;
ALTER TABLE products DROP COLUMN IF EXISTS uuid;
ALTER TABLE categories DROP COLUMN IF EXISTS uuid;

-- Расширение не удаляем, оно может другим таблицам нужно