package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port           string
	DatabasePath   string
	AlgoNodeURL    string
	AlgoIndexerURL string
	AlgoToken      string
	PaymentTimeout int // in minutes
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	timeoutStr := getEnv("PAYMENT_TIMEOUT", "30")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30 // default 30 minutes
	}
	// Hardcoded testnet configs 
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabasePath:   getEnv("DATABASE_PATH", "./algopay.db"),
		AlgoNodeURL:    getEnv("ALGO_NODE_URL", "https://testnet-api.algonode.cloud"),
		AlgoIndexerURL: getEnv("ALGO_INDEXER_URL", "https://testnet-idx.algonode.cloud"),
		AlgoToken:      getEnv("ALGO_TOKEN", ""),
		PaymentTimeout: timeout,
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
