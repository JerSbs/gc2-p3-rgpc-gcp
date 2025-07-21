package infra

import (
	"log"
	"sync"

	"gateway-service/config"
	"gateway-service/proto"

	"google.golang.org/grpc"
)

var (
	authClient    proto.AuthServiceClient
	paymentClient proto.PaymentServiceClient
	once          sync.Once
)

// InitGRPCClients initializes the gRPC clients for auth and payment services
func InitGRPCClients() {
	once.Do(func() {
		// Connect to Auth Service
		authConn, err := grpc.Dial(config.AuthServiceURL, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect to auth-service: %v", err)
		}
		authClient = proto.NewAuthServiceClient(authConn)
		log.Println("Connected to auth-service")

		// Connect to Payment Service
		// paymentConn, err := grpc.Dial(config.PaymentServiceURL, grpc.WithInsecure())
		paymentConn, err := grpc.Dial(config.PaymentServiceGRPC, grpc.WithInsecure())

		if err != nil {
			log.Fatalf("failed to connect to payment-service: %v", err)
		}
		paymentClient = proto.NewPaymentServiceClient(paymentConn)
		log.Println("Connected to payment-service")
	})
}

// GetAuthClient returns the initialized auth gRPC client
func GetAuthClient() proto.AuthServiceClient {
	return authClient
}

// GetPaymentClient returns the initialized payment gRPC client
func GetPaymentClient() proto.PaymentServiceClient {
	return paymentClient
}
