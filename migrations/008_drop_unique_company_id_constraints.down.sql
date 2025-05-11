-- Drop the partial unique index
DROP INDEX IF EXISTS unique_company_id_not_null;

-- Restore strict UNIQUE constraint on company_id (use with caution!)
ALTER TABLE users ADD CONSTRAINT unique_company_id UNIQUE (company_id);
