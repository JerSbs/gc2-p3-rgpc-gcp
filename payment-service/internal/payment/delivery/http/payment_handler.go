package http

import (
	"net/http"
	"payment-service/internal/payment/app"
	"payment-service/internal/payment/domain"
	"payment-service/internal/payment/infra"

	"github.com/labstack/echo/v4"
)

// Inisialisasi service payment
var paymentService = app.NewPaymentService(infra.NewPaymentRepository())

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a new payment and store it in the database
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body domain.Payment true "Payment Payload"
// @Success 201 {object} domain.Payment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments [post]
func CreatePayment(c echo.Context) error {
	var payload domain.Payment

	// Ambil request body dan bind ke struct
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": app.ErrInvalidPayload.Error()})
	}

	// Panggil service untuk validasi dan simpan ke DB
	result, err := paymentService.CreatePayment(c.Request().Context(), payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	// Kirim response sukses
	return c.JSON(http.StatusCreated, result)
}
