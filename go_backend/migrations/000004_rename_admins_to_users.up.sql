ALTER TABLE admins RENAME TO users;
ALTER TABLE urls RENAME COLUMN admin_id TO user_id;
ALTER INDEX IF EXISTS idx_urls_admin_id RENAME TO idx_urls_user_id;