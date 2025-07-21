package grpc

import (
	"fmt"
	"log"
	"net"
	"os"

	"payment-service/internal/payment/app"
	"payment-service/internal/payment/delivery/grpc/paymentpb"
	"payment-service/internal/payment/infra"

	"google.golang.org/grpc"
)

// Jalankan gRPC server
func StartGRPCServer() {
	// Ambil port dari env
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// Buka koneksi TCP
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Inisialisasi gRPC server
	grpcServer := grpc.NewServer()

	// Setup service dan handler
	repo := infra.NewPaymentRepository()
	service := app.NewPaymentService(repo)
	handler := &PaymentHandler{Service: service}

	// Daftarkan handler ke gRPC
	paymentpb.RegisterPaymentServiceServer(grpcServer, handler)

	fmt.Println("gRPC server running at port:", grpcPort)

	// Mulai server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
