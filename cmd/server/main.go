package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/database"
	"github.com/nomenarkt/lamina/config"
	"github.com/nomenarkt/lamina/internal/admin"
	"github.com/nomenarkt/lamina/internal/auth"
	"github.com/nomenarkt/lamina/internal/user"
)

func main() {
	fmt.Println("JWT_SECRET (startup):", os.Getenv("JWT_SECRET"))

	config.LoadEnv()

	db := database.ConnectDB()
	defer db.Close()

	router := gin.Default()

	api := router.Group("/api/v1")

	// Public routes
	auth.RegisterRoutes(api, db)

	// Protected routes
	api.Use(auth.AuthMiddleware())
	{
		user.RegisterRoutes(api, db)
		admin.RegisterRoutes(api, db)
		// tenant.RegisterRoutes(api, db) (future)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	router.Run(":" + port)
}
