-- 004_create_crew_assignments.sql

CREATE TABLE IF NOT EXISTS crew_assignments (
    id SERIAL PRIMARY KEY,
    flight_id INTEGER REFERENCES flights(id) ON DELETE CASCADE,
    crew_id INTEGER REFERENCES users(company_id) ON DELETE CASCADE,
    crew_role VARCHAR(10) NOT NULL,           -- e.g., CDB, OPL, CCA, PNC
    in_function BOOLEAN DEFAULT TRUE,         -- false = deadheading / MEP
    pickup_time TIMESTAMP,
    checkin_time TIMESTAMP,
    checkout_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT now()
);
