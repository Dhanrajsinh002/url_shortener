ALTER TABLE urls ADD COLUMN IF NOT EXISTS admin_id BIGINT REFERENCES admins(id);
CREATE INDEX IF NOT EXISTS idx_urls_admin_id ON urls(admin_id);