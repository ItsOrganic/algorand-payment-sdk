package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"algopay/algorand"
	"algopay/config"
	"algopay/db"
	"algopay/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// Server holds the API server dependencies
type Server struct {
	database    *db.Database
	algoClient  *algorand.Client
	config      *config.Config
	paymentChan chan *models.Payment
}

// NewServer creates a new API server
func NewServer(database *db.Database, algoClient *algorand.Client, config *config.Config) *Server {
	server := &Server{
		database:    database,
		algoClient:  algoClient,
		config:      config,
		paymentChan: make(chan *models.Payment, 100),
	}

	// Start webhook processor
	go server.processWebhooks()

	// Start payment monitor
	go algoClient.StartPaymentMonitor(server.paymentChan, database)

	// Start cleanup routine
	go server.cleanupExpiredPayments()

	return server
}

// SetupRoutes sets up the API routes
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/init-payment", s.initPayment)
		api.GET("/check-payment/:id", s.checkPayment)
		api.GET("/payment/:id", s.getPayment)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

// initPayment handles payment initialization
func (s *Server) initPayment(c *gin.Context) {
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// When api keys are there uncomment below piece to verify the address
	
	// Validate merchant address
	// if err := s.algoClient.ValidateAddress(req.MerchantAddress); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant address"})
	// 	return
	// }

	// Validate amount
	if req.Amount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	// Create payment record
	payment := &models.Payment{
		ID:              uuid.New().String(),
		MerchantAddress: req.MerchantAddress,
		Amount:          req.Amount,
		AssetID:         req.AssetID,
		CallbackURL:     req.CallbackURL,
		Status:          models.PaymentStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(time.Duration(s.config.PaymentTimeout) * time.Minute),
	}

	// Save to database
	if err := s.database.CreatePayment(payment); err != nil {
		log.Printf("Error creating payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// Generate QR code
	qrData := map[string]interface{}{
		"payment_id": payment.ID,
		"address":    payment.MerchantAddress,
		"amount":     payment.Amount,
		"asset_id":   payment.AssetID,
	}
	qrJSON, _ := json.Marshal(qrData)
	qrCode, err := qrcode.Encode(string(qrJSON), qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
	}

	response := models.PaymentResponse{
		PaymentID:       payment.ID,
		MerchantAddress: payment.MerchantAddress,
		Amount:          payment.Amount,
		AssetID:         payment.AssetID,
		ExpiresAt:       payment.ExpiresAt.Format(time.RFC3339),
		Status:          string(payment.Status),
	}

	if err == nil {
		response.QRCode = base64.StdEncoding.EncodeToString(qrCode)
	}

	c.JSON(http.StatusCreated, response)
}

// checkPayment handles payment status checking
func (s *Server) checkPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	payment, err := s.database.GetPayment(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	response := models.PaymentStatusResponse{
		PaymentID: payment.ID,
		Status:    payment.Status,
		TxnID:     payment.TxnID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// getPayment handles full payment details retrieval
func (s *Server) getPayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	payment, err := s.database.GetPayment(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// processWebhooks processes webhook notifications
func (s *Server) processWebhooks() {
	for payment := range s.paymentChan {
		// Update payment status in database
		if err := s.database.UpdatePaymentStatus(payment.ID, payment.Status, payment.TxnID); err != nil {
			log.Printf("Error updating payment status: %v", err)
			continue
		}

		// Send webhook if callback URL is provided
		if payment.CallbackURL != "" {
			go s.sendWebhook(payment)
		}
	}
}

// sendWebhook sends a webhook notification
func (s *Server) sendWebhook(payment *models.Payment) {
	webhook := models.WebhookPayload{
		PaymentID:       payment.ID,
		Status:          payment.Status,
		MerchantAddress: payment.MerchantAddress,
		Amount:          payment.Amount,
		AssetID:         payment.AssetID,
		TxnID:           payment.TxnID,
		Timestamp:       time.Now(),
	}

	jsonData, err := json.Marshal(webhook)
	if err != nil {
		log.Printf("Error marshaling webhook data: %v", err)
		return
	}

	// Send HTTP POST request
	resp, err := http.Post(payment.CallbackURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending webhook to %s: %v", payment.CallbackURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Webhook sent successfully for payment %s", payment.ID)
	} else {
		log.Printf("Webhook failed for payment %s, status: %d", payment.ID, resp.StatusCode)
	}
}

// cleanupExpiredPayments runs a cleanup routine for expired payments
func (s *Server) cleanupExpiredPayments() {
	ticker := time.NewTicker(5 * time.Minute) // Run every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.database.ExpireOldPayments(); err != nil {
				log.Printf("Error expiring old payments: %v", err)
			}
		}
	}
}
