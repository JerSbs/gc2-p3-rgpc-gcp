package http

import "github.com/labstack/echo/v4"

// response error standar
func ErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]any{
		"error": map[string]any{
			"code":    status,
			"message": message,
		},
	})
}
