package main

import (
	"log"
	"os"

	"github.com/nomenarkt/lamina/common/database"
	"github.com/nomenarkt/lamina/config"
	"github.com/nomenarkt/lamina/internal/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db := database.ConnectDB()
	defer db.Close()

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		auth.RegisterRoutes(api, db)
		// user.RegisterRoutes(api, db)
		// tenant.RegisterRoutes(api, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	router.Run(":" + port)
}
