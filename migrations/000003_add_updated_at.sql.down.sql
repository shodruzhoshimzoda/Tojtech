DROP TRIGGER IF EXISTS set_timestamp_categories ON categories;
DROP FUNCTION IF EXISTS trigger_set_timestamp;
ALTER TABLE categories DROP COLUMN IF EXISTS updated_at;