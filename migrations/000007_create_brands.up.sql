-- Brands: HP, ASUS, MSI ... A brand is NOT tied to one category:
-- ASUS makes graphics cards, motherboards and laptops at the same time.
CREATE TABLE brands (
    id          SERIAL PRIMARY KEY,
    uuid        UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,   -- used in filter urls: /products?brand=asus
    logo_url    TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- "HP" and "hp" must be the same brand, so the name is unique ignoring case
CREATE UNIQUE INDEX idx_brands_name_lower ON brands (LOWER(name));

-- reuse the function that was created in migration 000003
CREATE TRIGGER set_timestamp_brands
    BEFORE UPDATE ON brands
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- A product may have a brand (optional: old products stay without one).
-- ON DELETE SET NULL: if a brand is deleted, its products stay, just without a brand.
ALTER TABLE products
    ADD COLUMN brand_id INTEGER REFERENCES brands(id) ON DELETE SET NULL;

CREATE INDEX idx_products_brand ON products(brand_id);