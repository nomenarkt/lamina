INSERT INTO users (id, email, password_hash, role, status, created_at)
VALUES (
    3190,
    'm.rakotoarison@madagascarairlines.com',
    '$2a$12$n8n79q/toXgC4d.kkdY8NOPCiUOJZgBUkNG831Ynq1O0m61dHdiu6',
    'admin',
    'active',
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

-- Sync the sequence with the inserted ID
SELECT setval('public.users_id_seq', 3190);
