package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database // Variabel global untuk akses database

// ConnectDB digunakan untuk koneksi ke MongoDB
func ConnectDB() {
	// Membuat opsi koneksi dari MONGOURI (diambil dari .env)
	clientOptions := options.Client().ApplyURI(GetEnv("MONGOURI"))

	// Menentukan batas waktu koneksi 10 detik
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Melakukan koneksi ke MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Gagal konek ke MongoDB:", err)
	}

	// Mengatur database sesuai nama dari ENV
	DB = client.Database(GetEnv("MONGODB_NAME"))
}
