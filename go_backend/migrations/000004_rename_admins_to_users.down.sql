ALTER INDEX IF EXISTS idx_urls_user_id RENAME TO idx_urls_admin_id;
ALTER TABLE urls RENAME COLUMN user_id TO admin_id;
ALTER TABLE users RENAME TO admins;