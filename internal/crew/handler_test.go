package crew_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/crew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== MOCK SETUP =====

type MockCrewService struct {
	mock.Mock
}

func (m *MockCrewService) GetCrewByFlight(ctx context.Context, flightID int64) ([]crew.CrewAssignment, error) {
	args := m.Called(ctx, flightID)
	return args.Get(0).([]crew.CrewAssignment), args.Error(1)
}

func (m *MockCrewService) AssignCrew(ctx context.Context, ca *crew.CrewAssignment) error {
	args := m.Called(ctx, ca)
	return args.Error(0)
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

// ===== TEST CASE =====

func TestGetCrewByFlight_Success(t *testing.T) {
	mockService := new(MockCrewService)
	handler := crew.NewHandler(mockService)
	router := setupRouterWithHandler(handler)

	// Only mock what's actually used
	mockService.On("GetDetailedCrewByFlight", mock.Anything, int64(42)).Return([]crew.CrewAssignmentDetail{
		{
			CrewID:        1001,
			CrewRole:      "CDB",
			FlightNumber:  "MD710",
			DepartureCode: "TNR",
			ArrivalCode:   "FTU",
		},
	}, nil)

	mockService.On("GetCrewByFlight", mock.Anything, int64(42)).Return([]crew.CrewAssignment{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/crew/flight/42", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("RESPONSE BODY: %s", w.Body.String()) // 🧪 Debug
	assert.Equal(t, http.StatusOK, w.Code)

	expected := `[{
		"id": 0,
		"crew_id": 1001,
		"crew_role": "CDB",
		"in_function": false,
		"pickup_time": "0001-01-01T00:00:00Z",
		"checkin_time": "0001-01-01T00:00:00Z",
		"checkout_time": "0001-01-01T00:00:00Z",
		"flight_number": "MD710",
		"departure_code": "TNR",
		"arrival_code": "FTU",
		"crew_name": "",
		"crew_email": ""
	}]`

	assert.JSONEq(t, expected, w.Body.String())
}
