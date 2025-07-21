package main

import (
	"fmt"
	"gateway-service/config"
	http "gateway-service/internal/gateway/delivery/http"
	"gateway-service/internal/gateway/infra"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Gateway Service API
// @version 1.0
// @description REST API Gateway untuk Auth, Payment, dan Shopping Services
// @host localhost:8082
// @BasePath /
// @schemes http
func main() {
	// Load .env config
	config.ConnectConfig()

	// Init gRPC clients (auth, payment)
	infra.InitGRPCClients()

	// Init Echo
	e := echo.New()

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register REST routes
	http.RegisterGatewayRoutes(e)

	// Start server
	startServer(e, config.AppPort)
}

// Jalankan server di port dari .env
func startServer(e *echo.Echo, port string) {
	fmt.Printf("Starting Gateway on %s\n", port)
	e.Logger.Fatal(e.Start(port)) // port format dari .env: ":8082"
}
