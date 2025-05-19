-- 📦 Alter users table to support RBAC and profile completeness

ALTER TABLE users
ADD COLUMN IF NOT EXISTS user_type TEXT DEFAULT 'external' NOT NULL, -- admin, internal, external
ADD COLUMN IF NOT EXISTS phone TEXT,
ADD COLUMN IF NOT EXISTS address TEXT;
