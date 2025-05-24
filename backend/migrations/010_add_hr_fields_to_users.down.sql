ALTER TABLE users
  DROP COLUMN IF EXISTS access_expires_at,
  DROP COLUMN IF EXISTS profile_picture_url,
  DROP COLUMN IF EXISTS sex,
  DROP COLUMN IF EXISTS birthday,
  DROP COLUMN IF EXISTS marital_status,
  DROP COLUMN IF EXISTS spouse_name,
  DROP COLUMN IF EXISTS has_children,
  DROP COLUMN IF EXISTS number_of_children,
  DROP COLUMN IF EXISTS national_id,
  DROP COLUMN IF EXISTS emergency_contact_name,
  DROP COLUMN IF EXISTS emergency_contact_phone;
