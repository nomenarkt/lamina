package crew_test

import (
	"context"
	"testing"

	"github.com/nomenarkt/lamina/internal/crew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCrewRepo implements the crew.Repository interface
type MockCrewRepo struct {
	mock.Mock
}

func (m *MockCrewRepo) Create(ctx context.Context, assignment *crew.CrewAssignment) error {
	args := m.Called(ctx, assignment)
	return args.Error(0)
}

func (m *MockCrewRepo) GetByFlightID(ctx context.Context, flightID int64) ([]crew.CrewAssignment, error) {
	args := m.Called(ctx, flightID)
	return args.Get(0).([]crew.CrewAssignment), args.Error(1)
}

func (m *MockCrewRepo) DeleteByFlightID(ctx context.Context, flightID int64) error {
	args := m.Called(ctx, flightID)
	return args.Error(0)
}

func (m *MockCrewRepo) GetFlightIDByNumber(ctx context.Context, flightNumber string) (int64, error) {
	args := m.Called(ctx, flightNumber)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCrewRepo) GetDetailedByFlightID(ctx context.Context, flightID int64) ([]crew.CrewAssignmentDetail, error) {
	args := m.Called(ctx, flightID)
	return args.Get(0).([]crew.CrewAssignmentDetail), args.Error(1)
}

func TestService_AssignCrew_Success(t *testing.T) {
	repo := new(MockCrewRepo)
	service := crew.NewService(repo)

	assignment := &crew.CrewAssignment{
		FlightID:   1001,
		CrewID:     3199,
		CrewRole:   "CDB",
		InFunction: true,
	}

	repo.On("Create", mock.Anything, assignment).Return(nil)

	err := service.AssignCrew(context.Background(), assignment)

	assert.NoError(t, err)
	repo.AssertCalled(t, "Create", mock.Anything, assignment)
}

func TestService_GetCrewByFlight_Success(t *testing.T) {
	repo := new(MockCrewRepo)
	service := crew.NewService(repo)

	expected := []crew.CrewAssignment{
		{ID: 1, FlightID: 1001, CrewID: 3199, CrewRole: "CDB"},
		{ID: 2, FlightID: 1001, CrewID: 3200, CrewRole: "OPL"},
	}

	repo.On("GetByFlightID", mock.Anything, int64(1001)).Return(expected, nil)

	result, err := service.GetCrewByFlight(context.Background(), 1001)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestService_RemoveCrewByFlight_Success(t *testing.T) {
	repo := new(MockCrewRepo)
	service := crew.NewService(repo)

	repo.On("DeleteByFlightID", mock.Anything, int64(1001)).Return(nil)

	err := service.RemoveCrewByFlight(context.Background(), 1001)

	assert.NoError(t, err)
	repo.AssertCalled(t, "DeleteByFlightID", mock.Anything, int64(1001))
}

func TestService_ResolveFlightID_Success(t *testing.T) {
	repo := new(MockCrewRepo)
	service := crew.NewService(repo)

	repo.On("GetFlightIDByNumber", mock.Anything, "MD710").Return(int64(1), nil)

	flightID, err := service.ResolveFlightID(context.Background(), "MD710")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), flightID)
	repo.AssertCalled(t, "GetFlightIDByNumber", mock.Anything, "MD710")
}

func TestService_ResolveFlightID_NotFound(t *testing.T) {
	repo := new(MockCrewRepo)
	service := crew.NewService(repo)

	repo.On("GetFlightIDByNumber", mock.Anything, "UNKNOWN").Return(int64(0), assert.AnError)

	flightID, err := service.ResolveFlightID(context.Background(), "UNKNOWN")

	assert.Error(t, err)
	assert.Equal(t, int64(0), flightID)
	repo.AssertCalled(t, "GetFlightIDByNumber", mock.Anything, "UNKNOWN")
}
