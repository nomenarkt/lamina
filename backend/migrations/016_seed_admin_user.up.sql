INSERT INTO users (id, company_id, email, password_hash, role, status, created_at, user_type)
VALUES (
    1,  -- internal system ID
    3190,  -- employee_ID
    'm.rakotoarison@madagascarairlines.com',
    '$2a$12$n8n79q/toXgC4d.kkdY8NOPCiUOJZgBUkNG831Ynq1O0m61dHdiu6',
    'admin',
    'active',
    CURRENT_TIMESTAMP,
    'admin' -- user_type
)
ON CONFLICT (email) DO NOTHING;

-- Sync internal auto-incrementing sequence to avoid collisions
SELECT setval('public.users_id_seq', GREATEST((SELECT MAX(id) FROM users), 1));
