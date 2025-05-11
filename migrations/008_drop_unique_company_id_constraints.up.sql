-- Drop existing unique constraints (if still present)
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_company_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_company_id_key;

-- Add partial unique index: only enforce uniqueness when company_id is NOT NULL
CREATE UNIQUE INDEX IF NOT EXISTS unique_company_id_not_null
ON users(company_id)
WHERE company_id IS NOT NULL;
