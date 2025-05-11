package crew

import "context"

type CrewServiceInterface interface {
	AssignCrew(ctx context.Context, assignment *CrewAssignment) error
	GetCrewByFlight(ctx context.Context, flightID int64) ([]CrewAssignment, error)
	RemoveCrewByFlight(ctx context.Context, flightID int64) error
	ResolveFlightID(ctx context.Context, flightNumber string) (int64, error)
	GetDetailedCrewByFlight(ctx context.Context, flightID int64) ([]CrewAssignmentDetail, error)
}

type crewService struct {
	repo Repository
}

func NewService(repo Repository) CrewServiceInterface {
	return &crewService{repo: repo}
}

func (s *crewService) AssignCrew(ctx context.Context, assignment *CrewAssignment) error {
	return s.repo.Create(ctx, assignment)
}

func (s *crewService) GetCrewByFlight(ctx context.Context, flightID int64) ([]CrewAssignment, error) {
	return s.repo.GetByFlightID(ctx, flightID)
}

func (s *crewService) RemoveCrewByFlight(ctx context.Context, flightID int64) error {
	return s.repo.DeleteByFlightID(ctx, flightID)
}

func (s *crewService) ResolveFlightID(ctx context.Context, flightNumber string) (int64, error) {
	return s.repo.GetFlightIDByNumber(ctx, flightNumber)
}

func (s *crewService) GetDetailedCrewByFlight(ctx context.Context, flightID int64) ([]CrewAssignmentDetail, error) {
	return s.repo.GetDetailedByFlightID(ctx, flightID)
}
