package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"algopay/algorand"
	"algopay/api"
	"algopay/config"
	"algopay/db"

	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	var (
		envFile = flag.String("env", ".env", "Path to environment file")
		port    = flag.String("port", "", "Port to run the server on (overrides env)")
	)
	flag.Parse()

	// Load environment file if it exists
	if _, err := os.Stat(*envFile); err == nil {
		if err := godotenv.Load(*envFile); err != nil {
			log.Printf("Warning: Error loading %s file: %v", *envFile, err)
		}
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Override port if provided via command line
	if *port != "" {
		cfg.Port = *port
	}

	// Initialize database
	database, err := db.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize Algorand client
	algoClient, err := algorand.NewClient(cfg.AlgoNodeURL, cfg.AlgoIndexerURL, cfg.AlgoToken)
	if err != nil {
		log.Fatalf("Failed to initialize Algorand client: %v", err)
	}

	// Create API server
	server := api.NewServer(database, algoClient, cfg)

	// Setup routes
	router := server.SetupRoutes()

	// Print startup information
	fmt.Printf("🚀 AlgoPay Server Starting\n")
	fmt.Printf("📡 Port: %s\n", cfg.Port)
	fmt.Printf("🗄️ Database: %s\n", cfg.DatabasePath)
	fmt.Printf("🌐 Algorand Node: %s\n", cfg.AlgoNodeURL)
	fmt.Printf("🔍 Algorand Indexer: %s\n", cfg.AlgoIndexerURL)
	fmt.Printf("⏰ Payment Timeout: %d minutes\n", cfg.PaymentTimeout)
	fmt.Printf("\n📋 API Endpoints:\n")
	fmt.Printf("   POST /api/v1/init-payment     - Initialize new payment\n")
	fmt.Printf("   GET  /api/v1/check-payment/:id - Check payment status\n")
	fmt.Printf("   GET  /api/v1/payment/:id       - Get payment details\n")
	fmt.Printf("   GET  /health                   - Health check\n")
	fmt.Printf("\n🌟 Server running on http://localhost:%s\n", cfg.Port)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}