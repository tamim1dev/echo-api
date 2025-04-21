// middleware/jwt_middleware.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware creates a middleware to validate JWT tokens
func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header missing",
				})
			}

			// Check if the header is in the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization format, expected 'Bearer TOKEN'",
				})
			}

			tokenString := parts[1]

			// Parse the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Check if the token is valid
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Add claims to context for use in handlers
				c.Set("user", claims)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
		}
	}
}

// GetUserID extracts the user ID from the context
func GetUserID(c echo.Context) (uint, error) {
	user := c.Get("user")
	if user == nil {
		return 0, fmt.Errorf("user not found in context")
	}

	claims := user.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	return userID, nil
}
