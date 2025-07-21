package main

import (
	"fmt"
	"payment-service/config"
	"payment-service/internal/payment/app"
	grpcServer "payment-service/internal/payment/delivery/grpc"

	_ "payment-service/docs" // untuk Swagger
)

// @title Payment Service API
// @version 1.0
// @description REST API untuk pembayaran
// @host localhost:8081
// @BasePath /
func main() {
	// Load env dan koneksi MongoDB
	config.LoadEnv()
	config.ConnectDB()

	// Jalankan cron job pembersih data lama
	app.StartPaymentCleanupJob()

	// Jalankan gRPC di background
	go grpcServer.StartGRPCServer()

	// Blok utama agar aplikasi tidak langsung exit
	fmt.Println("Payment Service is running with gRPC and cron job...")

	select {} // blok selamanya

	// Init Echo
	// app := echo.New()
	// app.Use(middleware.Logger())
	// app.Use(middleware.Recover())

	// // Health Check
	// app.GET("/", func(c echo.Context) error {
	// 	return c.String(200, "Payment Service is running!")
	// })

	// // Swagger
	// app.GET("/swagger/*", echoSwagger.WrapHandler)

	// // Route utama (harusnya modular)
	// app.POST("/payments", handler.CreatePayment)

	// // Port
	// port := config.GetEnvOrDefault("PORT", "8081")
	// fmt.Println("Server running at PORT:", port)

	// // Start
	// app.Logger.Fatal(app.Start(":" + port))
}
