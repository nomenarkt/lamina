package crew

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	crewGroup := rg.Group("/crew")
	crewGroup.POST("/assign", h.AssignCrew)
	crewGroup.GET("/flight/:flight_id", h.GetCrewByFlight)
	crewGroup.DELETE("/flight/:flight_id", h.RemoveCrewByFlight)
	crewGroup.GET("/flight/:flight_id/details", h.GetCrewDetailsByFlight)
}
