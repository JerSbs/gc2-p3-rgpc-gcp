package http

import (
	"gateway-service/internal/gateway/delivery/http/middleware"

	"github.com/labstack/echo/v4"
)

// RegisterGatewayRoutes mendaftarkan semua route untuk gateway
func RegisterGatewayRoutes(e *echo.Echo) {
	handler := NewGatewayHandler()

	// Public Routes
	e.POST("/login", handler.LoginHandler)       // gRPC Login
	e.POST("/register", handler.RegisterHandler) // gRPC Register

	// with JWT

	// Shopping → /products
	products := e.Group("/products")
	products.Use(middleware.JWTMiddleware)
	products.GET("", handler.ShoppingProxy)
	products.POST("", handler.ShoppingProxy)

	// Shopping → /transactions
	transactions := e.Group("/transactions")
	transactions.Use(middleware.JWTMiddleware)
	transactions.GET("", handler.ShoppingProxy)
	transactions.POST("", handler.ShoppingProxy)

	// Payment → /payments
	payments := e.Group("/payments")
	payments.Use(middleware.JWTMiddleware)
	payments.GET("", handler.GetAllPaymentsHandler)
	payments.POST("", handler.CreatePaymentHandler)
}
