package crew

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service CrewServiceInterface // 🔧 Use the correct interface name
}

type AssignCrewRequest struct {
	FlightNumber string `json:"flight_number"`
	CrewID       int64  `json:"crew_id"`
	CrewRole     string `json:"crew_role"`
	InFunction   bool   `json:"in_function"`
	PickupTime   string `json:"pickup_time"`
	CheckinTime  string `json:"checkin_time"`
	CheckoutTime string `json:"checkout_time"`
}

func NewHandler(s CrewServiceInterface) *Handler {
	return &Handler{service: s}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// POST /crew/assign
func (h *Handler) AssignCrew(c *gin.Context) {
	var req AssignCrewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	flightID, err := h.service.ResolveFlightID(c.Request.Context(), req.FlightNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	ca := &CrewAssignment{
		FlightID:     flightID,
		CrewID:       int(req.CrewID),
		CrewRole:     req.CrewRole,
		InFunction:   req.InFunction,
		PickupTime:   parseTime(req.PickupTime),
		CheckinTime:  parseTime(req.CheckinTime),
		CheckoutTime: parseTime(req.CheckoutTime),
	}

	err = h.service.AssignCrew(c.Request.Context(), ca)
	if err != nil {
		log.Printf("❌ AssignCrew error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Crew assigned"})
}

// GET /crew/flight/:flight_id
func (h *Handler) GetCrewByFlight(c *gin.Context) {
	flightIDStr := c.Param("flight_id")
	flightID, err := strconv.ParseInt(flightIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	assignments, err := h.service.GetCrewByFlight(c.Request.Context(), flightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch crew"})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// DELETE /crew/flight/:flight_id
func (h *Handler) RemoveCrewByFlight(c *gin.Context) {
	flightIDStr := c.Param("flight_id")
	flightID, err := strconv.ParseInt(flightIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	err = h.service.RemoveCrewByFlight(c.Request.Context(), flightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove crew"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Crew unassigned from flight"})
}

// GET /crew/flight/:flight_id/details
func (h *Handler) GetCrewDetailsByFlight(c *gin.Context) {
	flightIDStr := c.Param("flight_id")
	flightID, err := strconv.ParseInt(flightIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight ID"})
		return
	}

	result, err := h.service.GetDetailedCrewByFlight(c.Request.Context(), flightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve crew details"})
		return
	}

	c.JSON(http.StatusOK, result)
}
