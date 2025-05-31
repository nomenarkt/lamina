CREATE TABLE ranks (
    id SERIAL PRIMARY KEY,
    function_id INTEGER NOT NULL REFERENCES functions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    level INTEGER,
    UNIQUE(function_id, name)
);
