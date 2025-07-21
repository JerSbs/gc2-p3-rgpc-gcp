package config

import "fmt"

// Variabel global config
var (
	AuthServiceURL     string
	PaymentServiceURL  string
	ShoppingServiceURL string
	JWTSecret          string
	AppPort            string
	PaymentServiceGRPC string // untuk REST forward

)

// ConnectConfig akan load semua config dari file .env atau fallback ke default
func ConnectConfig() {
	LoadEnv()

	AuthServiceURL = GetEnvOrDefault("AUTH_SERVICE_ADDR", "localhost:50052")
	PaymentServiceGRPC = GetEnvOrDefault("PAYMENT_SERVICE_GRPC_ADDR", "localhost:50051")
	PaymentServiceURL = GetEnvOrDefault("PAYMENT_SERVICE_REST_URL", "http://localhost:8081")

	ShoppingServiceURL = GetEnvOrDefault("SHOPPING_SERVICE_URL", "http://localhost:8080")
	JWTSecret = GetEnvOrDefault("JWT_SECRET", "mysecretkey123")

	AppPort = GetEnvOrDefault("GATEWAY_PORT", ":8082")

	fmt.Println("Gateway config loaded")
}
