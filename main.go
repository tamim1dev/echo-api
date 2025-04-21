// main.go
package main

import (
	"log"
	"net/http"
	"os"

	"echo-api/config"
	"echo-api/controllers"
	"echo-api/middleware"
	"echo-api/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = models.MigrateDB(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Load JWT configuration
	jwtConfig, err := config.LoadJWTConfig()
	if err != nil {
		log.Fatalf("Failed to load JWT configuration: %v", err)
	}

	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Register controllers
	userController := controllers.NewUserController(db)
	authController := controllers.NewAuthController(db, jwtConfig)

	// Public API Routes
	api := e.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)

		// Protected routes
		users := api.Group("/users")
		users.Use(middleware.JWTMiddleware(jwtConfig.SecretKey))
		users.GET("", userController.GetUsers)
		users.GET("/:id", userController.GetUser)
		users.POST("", userController.CreateUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "UP",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
