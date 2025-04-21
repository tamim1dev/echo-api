// config/jwt.go
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds the JWT configuration
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// LoadJWTConfig loads JWT configuration from environment variables
func LoadJWTConfig() (*JWTConfig, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}

	// Default token duration is 24 hours
	tokenDuration := 24 * time.Hour

	return &JWTConfig{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}, nil
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID uint, email string, config *JWTConfig) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(config.TokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.SecretKey))
}
