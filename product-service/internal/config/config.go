package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort             string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	KafkaBrokers         []string
	KafkaTopic           string
	InventoryServiceAddr string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using env vars")
	}
	return &Config{
		GRPCPort:             getEnv("GRPC_PORT", "50051"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "productdb"),
		KafkaBrokers:         strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:           getEnv("KAFKA_TOPIC", "product-events"),
		InventoryServiceAddr: getEnv("INVENTORY-SERVICE-ADDR", "localhost:50052"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
