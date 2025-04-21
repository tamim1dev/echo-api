// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"echo-api/config"
	"echo-api/models"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	db        *gorm.DB
	jwtConfig *config.JWTConfig
}

// NewAuthController creates a new auth controller
func NewAuthController(db *gorm.DB, jwtConfig *config.JWTConfig) *AuthController {
	return &AuthController{
		db:        db,
		jwtConfig: jwtConfig,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (ac *AuthController) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Find user by email
	var user models.User
	result := ac.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := config.GenerateToken(user.ID, user.Email, ac.jwtConfig)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	// Return token and user data
	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Register creates a new user account
func (ac *AuthController) Register(c echo.Context) error {
	req := new(models.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Check if email already exists
	var existingUser models.User
	result := ac.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process user data",
		})
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	result = ac.db.Create(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := config.GenerateToken(user.ID, user.Email, ac.jwtConfig)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	// Return token and user data
	return c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}
