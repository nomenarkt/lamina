CREATE TABLE user_organizational_units (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES organizational_units(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, unit_id)
);
