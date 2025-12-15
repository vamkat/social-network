
-- Drop trigger
DROP TRIGGER IF EXISTS update_file_updated_at ON files;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop table
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS file_variants;
