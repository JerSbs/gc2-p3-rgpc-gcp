package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv untuk membaca file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv untuk mengambil value wajib dari ENV
func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable '%s' tidak ditemukan atau kosong", key)
	}
	return val
}

// GetEnvOrDefault untuk ambil ENV, jika kosong pakai fallback
func GetEnvOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
