package middleware

import (
	"net/http"
	"strings"

	"gateway-service/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware untuk validasi token dan ekstrak email dari JWT
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing Authorization header",
			})
		}

		// Ambil token dari "Bearer <token>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan pakai metode HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.ErrUnauthorized
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid or expired token",
			})
		}

		// Ambil data dari claims (email, user_id, dll)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if email, ok := claims["email"].(string); ok {
				c.Set("userEmail", email)
			}
		}

		return next(c)
	}
}
