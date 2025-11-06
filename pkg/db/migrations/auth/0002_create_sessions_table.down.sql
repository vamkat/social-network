-- 0002_create_sessions_table.down.sql
-- Rollback: drop sessions table

DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_access_token;
DROP TABLE IF EXISTS sessions CASCADE;
