-- 1. Включаем расширение для генерации uuid. 1 раз на всю БД
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Добавляем колонку. NOT NULL + DEFAULT + UNIQUE сразу
-- Для старых строк DEFAULT сработает и всем поставит UUID
ALTER TABLE categories
    ADD COLUMN uuid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE;

ALTER TABLE products
    ADD COLUMN uuid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE;

ALTER TABLE product_images
    ADD COLUMN uuid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE;

-- 3. Бонус: сразу добавим индексы на uuid. Поиск по uuid будет быстрым
CREATE INDEX idx_categories_uuid ON categories(uuid);
CREATE INDEX idx_products_uuid ON products(uuid);
CREATE INDEX idx_product_images_uuid ON product_images(uuid);