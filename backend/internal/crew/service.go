// Package crew contains business logic and service interfaces for crew assignment operations.
package crew

import "context"

// ServiceInterface defines the operations for assigning and retrieving crew information.
type ServiceInterface interface {
	AssignCrew(ctx context.Context, assignment *Assignment) error
	GetCrewByFlight(ctx context.Context, flightID int64) ([]Assignment, error)
	RemoveCrewByFlight(ctx context.Context, flightID int64) error
	ResolveFlightID(ctx context.Context, flightNumber string) (int64, error)
	GetDetailedCrewByFlight(ctx context.Context, flightID int64) ([]AssignmentDetail, error)
}

// crewService implements the ServiceInterface using a data repository.
type crewService struct {
	repo Repository
}

// NewService creates a new instance of ServiceInterface using the provided repository.
func NewService(repo Repository) ServiceInterface {
	return &crewService{repo: repo}
}

// AssignCrew stores a new crew assignment in the database.
func (s *crewService) AssignCrew(ctx context.Context, assignment *Assignment) error {
	return s.repo.Create(ctx, assignment)
}

// GetCrewByFlight returns a list of crew assignments for a given flight ID.
func (s *crewService) GetCrewByFlight(ctx context.Context, flightID int64) ([]Assignment, error) {
	return s.repo.GetByFlightID(ctx, flightID)
}

// RemoveCrewByFlight removes all crew assignments associated with a given flight ID.
func (s *crewService) RemoveCrewByFlight(ctx context.Context, flightID int64) error {
	return s.repo.DeleteByFlightID(ctx, flightID)
}

// ResolveFlightID looks up the internal flight ID based on a flight number.
func (s *crewService) ResolveFlightID(ctx context.Context, flightNumber string) (int64, error) {
	return s.repo.GetFlightIDByNumber(ctx, flightNumber)
}

// GetDetailedCrewByFlight retrieves detailed crew information for a specific flight.
func (s *crewService) GetDetailedCrewByFlight(ctx context.Context, flightID int64) ([]AssignmentDetail, error) {
	return s.repo.GetDetailedByFlightID(ctx, flightID)
}
