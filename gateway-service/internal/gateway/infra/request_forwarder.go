package infra

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ForwardRequest mengirim request ke backend service, lalu kembalikan responsenya ke client
func ForwardRequest(c echo.Context, targetURL string) error {
	req, err := http.NewRequest(c.Request().Method, targetURL, c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create request"})
	}

	// Copy semua header dari original request
	for k, v := range c.Request().Header {
		req.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"message": "Service unreachable"})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
