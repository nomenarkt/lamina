-- 007_make_company_id_nullable.down.sql

-- Optional safety check:
-- DELETE FROM users WHERE company_id IS NULL;

ALTER TABLE users
ALTER COLUMN company_id SET NOT NULL;
