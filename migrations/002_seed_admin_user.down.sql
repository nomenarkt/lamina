DELETE FROM users WHERE id = 3190;
SELECT setval('users_id_seq', 1, false);
