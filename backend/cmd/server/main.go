// Package main initializes and starts the Lamina backend API server.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nomenarkt/lamina/common/database"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/config"
	"github.com/nomenarkt/lamina/internal/admin"
	"github.com/nomenarkt/lamina/internal/auth"
	"github.com/nomenarkt/lamina/internal/crew"
	"github.com/nomenarkt/lamina/internal/tasks"
	"github.com/nomenarkt/lamina/internal/user"
)

func main() {
	fmt.Println("JWT_SECRET (startup):", os.Getenv("JWT_SECRET"))

	config.LoadEnv()

	db := database.ConnectDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")

	// Auth setup
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewService(authRepo)
	auth.RegisterRoutes(api, db, authService)

	// Secure endpoints
	userRepo := user.NewUserRepository(db)
	api.Use(auth.Middleware(userRepo)) // ✅ Inject userRepo

	{
		userService := user.NewUserService(userRepo)
		userHandler := user.NewUserHandler(userService)
		user.RegisterRoutes(api, userHandler)

		tasks.StartUserCleanupTask(userRepo)

		adminRepo := admin.NewAdminRepository(db)
		hasher := &utils.BcryptHasher{}
		adminService := admin.NewAdminService(adminRepo, hasher)
		admin.RegisterRoutes(api, adminService)

		crewRepo := crew.NewRepository(db)
		crewService := crew.NewService(crewRepo)
		crewHandler := crew.NewHandler(crewService)
		crew.RegisterRoutes(api, crewHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
