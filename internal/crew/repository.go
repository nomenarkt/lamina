package crew

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, ca *CrewAssignment) error
	GetByFlightID(ctx context.Context, flightID int64) ([]CrewAssignment, error)
	DeleteByFlightID(ctx context.Context, flightID int64) error
	GetFlightIDByNumber(ctx context.Context, flightNumber string) (int64, error)
	GetDetailedByFlightID(ctx context.Context, flightID int64) ([]CrewAssignmentDetail, error)
}

type crewRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &crewRepository{db: db}
}

func (r *crewRepository) Create(ctx context.Context, ca *CrewAssignment) error {
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

func (r *crewRepository) GetByFlightID(ctx context.Context, flightID int64) ([]CrewAssignment, error) {
	var assignments []CrewAssignment
	err := r.db.SelectContext(ctx, &assignments, `
		SELECT * FROM crew_assignments
		WHERE flight_id = $1
		ORDER BY pickup_time ASC
	`, flightID)
	return assignments, err
}

func (r *crewRepository) DeleteByFlightID(ctx context.Context, flightID int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM crew_assignments
		WHERE flight_id = $1
	`, flightID)
	return err
}

func (r *crewRepository) GetFlightIDByNumber(ctx context.Context, flightNumber string) (int64, error) {
	var id int64
	err := r.db.GetContext(ctx, &id, `SELECT id FROM flights WHERE flight_number = $1`, flightNumber)
	return id, err
}

func (r *crewRepository) GetDetailedByFlightID(ctx context.Context, flightID int64) ([]CrewAssignmentDetail, error) {
	var result []CrewAssignmentDetail
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
