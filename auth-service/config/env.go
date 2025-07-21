package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv memuat .env file jika ada
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with system environment variables.")
	}
}
