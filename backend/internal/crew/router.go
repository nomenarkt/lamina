// Package crew defines route registration for crew assignment APIs.
package crew

import "github.com/gin-gonic/gin"

// RegisterRoutes sets up the HTTP endpoints for crew management under the given route group.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	crewGroup := rg.Group("/crew")
	crewGroup.POST("/assign", h.AssignCrew)
	crewGroup.GET("/flight/:flight_id", h.GetCrewByFlight)
	crewGroup.DELETE("/flight/:flight_id", h.RemoveCrewByFlight)
	crewGroup.GET("/flight/:flight_id/details", h.GetCrewDetailsByFlight)
}
