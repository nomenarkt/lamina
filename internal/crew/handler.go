// Package crew handles HTTP endpoints for crew assignment and lookup.
package crew

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler defines the HTTP handler for crew operations.
type Handler struct {
	service ServiceInterface
}

// AssignCrewRequest represents the expected JSON payload for crew assignment.
type AssignCrewRequest struct {
	FlightNumber string `json:"flight_number"`
	CrewID       int64  `json:"crew_id"`
	CrewRole     string `json:"crew_role"`
	InFunction   bool   `json:"in_function"`
	PickupTime   string `json:"pickup_time"`
	CheckinTime  string `json:"checkin_time"`
	CheckoutTime string `json:"checkout_time"`
}

// NewHandler creates a new Handler instance for the crew service.
func NewHandler(s ServiceInterface) *Handler {
	return &Handler{service: s}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// AssignCrew assigns a crew member to a flight.
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

	ca := &Assignment{
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

// GetCrewByFlight retrieves assigned crew for a flight.
// GET /crew/flight/:flight_id
func (h *Handler) GetCrewByFlight(c *gin.Context) {
	flightIDStr := c.Param("flight_id")
	flightID, err := strconv.ParseInt(flightIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight ID"})
		return
	}

	crew, err := h.service.GetDetailedCrewByFlight(c.Request.Context(), flightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve crew"})
		return
	}

	c.JSON(http.StatusOK, crew)
}

// RemoveCrewByFlight unassigns all crew from a flight.
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

// GetCrewDetailsByFlight returns detailed crew info for a flight.
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
