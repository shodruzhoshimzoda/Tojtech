-- migrations/000005_widen_product_image_url.down.sql
ALTER TABLE product_images
    ALTER COLUMN image_url TYPE VARCHAR(500);