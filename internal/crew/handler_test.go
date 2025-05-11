package crew_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/crew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCrewService struct {
	mock.Mock
}

func (m *MockCrewService) AssignCrew(ctx context.Context, ca *crew.CrewAssignment) error {
	args := m.Called(ctx, ca)
	return args.Error(0)
}

func (m *MockCrewService) GetCrewByFlight(ctx context.Context, flightID int64) ([]crew.CrewAssignment, error) {
	args := m.Called(ctx, flightID)
	return args.Get(0).([]crew.CrewAssignment), args.Error(1)
}

func (m *MockCrewService) RemoveCrewByFlight(ctx context.Context, flightID int64) error {
	args := m.Called(ctx, flightID)
	return args.Error(0)
}

func (m *MockCrewService) GetDetailedCrewByFlight(ctx context.Context, flightID int64) ([]crew.CrewAssignmentDetail, error) {
	args := m.Called(ctx, flightID)
	return args.Get(0).([]crew.CrewAssignmentDetail), args.Error(1)
}

func (m *MockCrewService) ResolveFlightID(ctx context.Context, flightNumber string) (int64, error) {
	args := m.Called(ctx, flightNumber)
	return args.Get(0).(int64), args.Error(1)
}

func setupRouterWithHandler(handler *crew.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	group := r.Group("/crew")
	group.POST("/assign", handler.AssignCrew)
	group.GET("/flight/:flight_id", handler.GetCrewByFlight)
	group.DELETE("/flight/:flight_id", handler.RemoveCrewByFlight)

	return r
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestAssignCrew_Success(t *testing.T) {
	mockService := new(MockCrewService)
	handler := crew.NewHandler(mockService)
	router := setupRouterWithHandler(handler)

	request := crew.AssignCrewRequest{
		FlightNumber: "MD710",
		CrewID:       1001,
		CrewRole:     "CDB",
		InFunction:   true,
		PickupTime:   "2025-05-08T12:30:00Z",
		CheckinTime:  "2025-05-08T13:00:00Z",
		CheckoutTime: "2025-05-08T15:00:00Z",
	}

	mockService.On("ResolveFlightID", mock.Anything, "MD710").Return(int64(1), nil)

	expected := &crew.CrewAssignment{
		FlightID:     1,
		CrewID:       int(request.CrewID),
		CrewRole:     request.CrewRole,
		InFunction:   request.InFunction,
		PickupTime:   parseTime(request.PickupTime),
		CheckinTime:  parseTime(request.CheckinTime),
		CheckoutTime: parseTime(request.CheckoutTime),
	}
	mockService.On("AssignCrew", mock.Anything, expected).Return(nil)

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/crew/assign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Crew assigned")
	mockService.AssertExpectations(t)
}

func TestAssignCrew_InvalidPayload(t *testing.T) {
	mockService := new(MockCrewService)
	handler := crew.NewHandler(mockService)
	router := setupRouterWithHandler(handler)

	req := httptest.NewRequest(http.MethodPost, "/crew/assign", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid payload")
}

func TestGetCrewByFlight_Success(t *testing.T) {
	mockService := new(MockCrewService)
	handler := crew.NewHandler(mockService)
	router := setupRouterWithHandler(handler)

	mockService.On("GetDetailedCrewByFlight", mock.Anything, int64(42)).Return([]crew.CrewAssignmentDetail{
		{CrewID: 1001, CrewRole: "CDB", FlightNumber: "MD710", DepartureCode: "TNR", ArrivalCode: "FTU"},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/crew/flight/42", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "CDB")
}

func TestRemoveCrewByFlight_Success(t *testing.T) {
	mockService := new(MockCrewService)
	handler := crew.NewHandler(mockService)
	router := setupRouterWithHandler(handler)

	mockService.On("RemoveCrewByFlight", mock.Anything, int64(42)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/crew/flight/42", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "unassigned")
}
