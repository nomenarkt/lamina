// Package main starts the application server.
package main

import (
	"fmt"
	"log"
	"os"

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

	// ✅ Properly check error on closing DB
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	router := gin.Default()
	api := router.Group("/api/v1")

	// Public routes
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewService(authRepo)
	auth.RegisterRoutes(api, db, authService)

	// Protected routes
	api.Use(auth.Middleware())
	{
		userRepo := user.NewUserRepository(db)
		userService := user.NewUserService(userRepo)
		userHandler := user.NewUserHandler(userService)
		user.RegisterRoutes(api, userHandler)

		// ✅ Start 24h cleanup background job
		tasks.StartUserCleanupTask(userRepo)

		adminRepo := admin.NewAdminRepository(db)
		hasher := &utils.BcryptHasher{}
		adminService := admin.NewAdminService(adminRepo, hasher)
		admin.RegisterRoutes(api, adminService)

		crewRepo := crew.NewRepository(db)
		crewService := crew.NewService(crewRepo)
		crewHandler := crew.NewHandler(crewService)
		crew.RegisterRoutes(api, crewHandler)

		// tenant.RegisterRoutes(api, db)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)

	// ✅ Properly check error when running router
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
