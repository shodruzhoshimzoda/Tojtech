-- migrations/000005_widen_product_image_url.up.sql
ALTER TABLE product_images
    ALTER COLUMN image_url TYPE TEXT;