-- 🔧 Drop dependent foreign key constraints first (if they exist)
ALTER TABLE crew_assignments DROP CONSTRAINT IF EXISTS crew_assignments_crew_id_fkey;

-- ✅ Drop existing unique constraints on users.company_id
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_company_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_company_id_key;

-- ✅ Add partial unique index on users.company_id (only when NOT NULL)
CREATE UNIQUE INDEX IF NOT EXISTS unique_company_id_not_null
ON users(company_id)
WHERE company_id IS NOT NULL;

-- 🔁 Re-add foreign key from crew_assignments.crew_id → users.id
ALTER TABLE crew_assignments
ADD CONSTRAINT crew_assignments_crew_id_fkey
FOREIGN KEY (crew_id) REFERENCES users(id) ON DELETE CASCADE;
