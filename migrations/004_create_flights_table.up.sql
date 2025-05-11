-- +migrate Up
CREATE TABLE IF NOT EXISTS flights (
    id SERIAL PRIMARY KEY,
    flight_number VARCHAR(10) NOT NULL,
    departure_code VARCHAR(5) NOT NULL,      -- IATA code (e.g., TNR)
    arrival_code VARCHAR(5) NOT NULL,
    scheduled_departure TIMESTAMP NOT NULL,
    scheduled_arrival TIMESTAMP NOT NULL,
    actual_departure TIMESTAMP,
    actual_arrival TIMESTAMP,
    delay_reason TEXT,
    planned_load INTEGER,
    actual_load INTEGER,
    airplane_immatriculation VARCHAR(20),
    created_at TIMESTAMP DEFAULT now()
);
