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
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo)
	auth.RegisterRoutes(api, db, authService)

	// Protected routes
	api.Use(auth.AuthMiddleware())
	{
		userRepo := user.NewUserRepository(db)
		userService := user.NewUserService(userRepo)
		userHandler := user.NewUserHandler(userService)
		user.RegisterRoutes(api, userHandler)

		adminRepo := admin.NewAdminRepository(db)
		hasher := &utils.BcryptHasher{}
		adminService := admin.NewAdminService(adminRepo, hasher)
		admin.RegisterRoutes(api, adminService)

		crewRepo := crew.NewRepository(db)
		crewService := crew.NewService(crewRepo) // not NewCrewService
		crewHandler := crew.NewHandler(crewService)
		crew.RegisterRoutes(api, crewHandler)

		// tenant.RegisterRoutes(api, db) (future)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	router.Run(":" + port)
}
