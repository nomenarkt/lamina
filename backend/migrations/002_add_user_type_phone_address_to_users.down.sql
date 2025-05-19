-- Remove user_type, phone, and address from users table
ALTER TABLE users DROP COLUMN IF EXISTS user_type;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users DROP COLUMN IF EXISTS address;
