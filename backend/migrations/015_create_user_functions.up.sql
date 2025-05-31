CREATE TABLE user_functions (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    function_id INTEGER NOT NULL REFERENCES functions(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES organizational_units(id) ON DELETE CASCADE,
    rank_id INTEGER REFERENCES ranks(id),
    PRIMARY KEY(user_id, function_id, unit_id)
);
