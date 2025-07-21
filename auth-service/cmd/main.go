// cmd/main.go
package main

import (
	"auth-service/config"
	"auth-service/internal/auth/app"
	authgrpc "auth-service/internal/auth/delivery/grpc"
	"auth-service/internal/auth/delivery/grpc/authpb"
	"auth-service/internal/auth/infra"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// Load ENV
	config.LoadEnv()

	// Connect to MongoDB
	mongoURI := os.Getenv("AUTH_MONGO_URI")
	dbName := os.Getenv("AUTH_DB_NAME")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to create mongo client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	db := client.Database(dbName)

	// Init dependencies
	repo := infra.NewMongoUserRepository(db)
	jwtManager := app.NewJWTManager()
	hasher := &app.BcryptHasher{}
	service := app.NewAuthService(repo, jwtManager, hasher)
	handler := authgrpc.NewAuthHandler(service)

	// Start gRPC server
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, handler)

	fmt.Printf("Auth Service running on port %s...\n", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
