ALTER TABLE users
  ADD COLUMN access_expires_at TIMESTAMPTZ,
  ADD COLUMN profile_picture_url TEXT,
  ADD COLUMN sex TEXT,
  ADD COLUMN birthday DATE,
  ADD COLUMN marital_status BOOLEAN,
  ADD COLUMN spouse_name TEXT,
  ADD COLUMN has_children BOOLEAN,
  ADD COLUMN number_of_children INTEGER,
  ADD COLUMN national_id TEXT,
  ADD COLUMN emergency_contact_name TEXT,
  ADD COLUMN emergency_contact_phone TEXT;
