// Package crew implements repository interfaces for managing crew assignments.
package crew

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for interacting with crew assignment storage.
type Repository interface {
	Create(ctx context.Context, ca *Assignment) error
	GetByFlightID(ctx context.Context, flightID int64) ([]Assignment, error)
	DeleteByFlightID(ctx context.Context, flightID int64) error
	GetFlightIDByNumber(ctx context.Context, flightNumber string) (int64, error)
	GetDetailedByFlightID(ctx context.Context, flightID int64) ([]AssignmentDetail, error)
}

// crewRepository is a concrete implementation of the Repository interface.
type crewRepository struct {
	db *sqlx.DB
}

// NewRepository returns a new instance of a crew Repository.
func NewRepository(db *sqlx.DB) Repository {
	return &crewRepository{db: db}
}

// Create inserts a new Assignment into the database.
func (r *crewRepository) Create(ctx context.Context, ca *Assignment) error {
	query := `
		INSERT INTO crew_assignments (flight_id, crew_id, crew_role, in_function, pickup_time, checkin_time, checkout_time)
		VALUES (:flight_id, :crew_id, :crew_role, :in_function, :pickup_time, :checkin_time, :checkout_time)
		RETURNING id
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, &ca.ID, ca)
}

// GetByFlightID returns all crew assignments for a specific flight.
func (r *crewRepository) GetByFlightID(ctx context.Context, flightID int64) ([]Assignment, error) {
	var assignments []Assignment
	err := r.db.SelectContext(ctx, &assignments, `
		SELECT * FROM crew_assignments
		WHERE flight_id = $1
		ORDER BY pickup_time ASC
	`, flightID)
	return assignments, err
}

// DeleteByFlightID removes all crew assignments for a specific flight.
func (r *crewRepository) DeleteByFlightID(ctx context.Context, flightID int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM crew_assignments
		WHERE flight_id = $1
	`, flightID)
	return err
}

// GetFlightIDByNumber returns the internal ID of a flight given its flight number.
func (r *crewRepository) GetFlightIDByNumber(ctx context.Context, flightNumber string) (int64, error) {
	var id int64
	err := r.db.GetContext(ctx, &id, `SELECT id FROM flights WHERE flight_number = $1`, flightNumber)
	return id, err
}

// GetDetailedByFlightID returns enriched crew assignment data for a flight.
func (r *crewRepository) GetDetailedByFlightID(ctx context.Context, flightID int64) ([]AssignmentDetail, error) {
	var result []AssignmentDetail
	query := `
		SELECT 
			ca.id, ca.crew_id, ca.crew_role, ca.in_function,
			ca.pickup_time, ca.checkin_time, ca.checkout_time,
			f.flight_number, f.departure_code, f.arrival_code,
			u.full_name AS crew_name, u.email AS crew_email
		FROM crew_assignments ca
		JOIN flights f ON ca.flight_id = f.id
		JOIN users u ON ca.crew_id = u.company_id
		WHERE ca.flight_id = $1
		ORDER BY ca.pickup_time ASC
	`

	err := r.db.SelectContext(ctx, &result, query, flightID)
	return result, err
}
