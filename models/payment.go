package models

import (
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusExpired   PaymentStatus = "expired"
)

// Payment represents a payment request
type Payment struct {
	ID              string        `json:"id" db:"id"`
	MerchantAddress string        `json:"merchant_address" db:"merchant_address"`
	Amount          uint64        `json:"amount" db:"amount"`
	AssetID         uint64        `json:"asset_id" db:"asset_id"`
	CallbackURL     string        `json:"callback_url" db:"callback_url"`
	Status          PaymentStatus `json:"status" db:"status"`
	TxnID           string        `json:"txn_id,omitempty" db:"txn_id"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	ExpiresAt       time.Time     `json:"expires_at" db:"expires_at"`
}

// PaymentRequest represents a payment initialization request
type PaymentRequest struct {
	MerchantAddress string `json:"merchant_address" binding:"required"`
	Amount          uint64 `json:"amount" binding:"required"`
	AssetID         uint64 `json:"asset_id"`
	CallbackURL     string `json:"callback_url"`
}

// PaymentResponse represents a payment initialization response
type PaymentResponse struct {
	PaymentID       string `json:"payment_id"`
	MerchantAddress string `json:"merchant_address"`
	Amount          uint64 `json:"amount"`
	AssetID         uint64 `json:"asset_id"`
	QRCode          string `json:"qr_code,omitempty"`
	ExpiresAt       string `json:"expires_at"`
	Status          string `json:"status"`
}

// PaymentStatusResponse represents a payment status check response
type PaymentStatusResponse struct {
	PaymentID string        `json:"payment_id"`
	Status    PaymentStatus `json:"status"`
	TxnID     string        `json:"txn_id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// WebhookPayload represents the payload sent to callback URLs
type WebhookPayload struct {
	PaymentID       string        `json:"payment_id"`
	Status          PaymentStatus `json:"status"`
	MerchantAddress string        `json:"merchant_address"`
	Amount          uint64        `json:"amount"`
	AssetID         uint64        `json:"asset_id"`
	TxnID           string        `json:"txn_id"`
	Timestamp       time.Time     `json:"timestamp"`
}