package algorand

import (
	"context"
	"fmt"
	"log"
	"time"

	"algopay/models"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Client wraps Algorand SDK clients
type Client struct {
	algodClient   *algod.Client
	indexerClient *indexer.Client
}

// Transaction represents a simplified transaction for our use case
type Transaction struct {
	ID        string
	Sender    string
	Receiver  string
	Amount    uint64
	AssetID   uint64
	Round     uint64
	Timestamp time.Time
}

// NewClient creates a new Algorand client
func NewClient(nodeURL, indexerURL, token string) (*Client, error) {
	// Create algod client
	algodClient, err := algod.MakeClient(nodeURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create algod client: %w", err)
	}

	// Create indexer client
	indexerClient, err := indexer.MakeClient(indexerURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexer client: %w", err)
	}

	return &Client{
		algodClient:   algodClient,
		indexerClient: indexerClient,
	}, nil
}

// GetAccountInfo gets account information
// func (c *Client) GetAccountInfo(address string) (types.Account, error) {
// 	accountInfo, err := c.algodClient.AccountInformation(address).Do(context.Background())
// 	if err != nil {
// 		return types.Account{}, fmt.Errorf("failed to get account info: %w", err)
// 	}
// 	return accountInfo, nil
// }

// CheckPayment checks if a payment has been made to the specified address
func (c *Client) CheckPayment(payment *models.Payment, lastCheckedRound uint64) (*Transaction, error) {
	// Query transactions for the merchant address
	result, err := c.indexerClient.LookupAccountTransactions(payment.MerchantAddress).
		MinRound(lastCheckedRound + 1).
		TxType("pay"). // For ALGO payments
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to lookup transactions: %w", err)
	}

	// Check ALGO payments (asset ID 0)
	if payment.AssetID == 0 {
		for _, txn := range result.Transactions {
			if txn.PaymentTransaction.Receiver == payment.MerchantAddress &&
				txn.PaymentTransaction.Amount >= payment.Amount {
				return &Transaction{
					ID:        txn.Id,
					Sender:    txn.Sender,
					Receiver:  txn.PaymentTransaction.Receiver,
					Amount:    txn.PaymentTransaction.Amount,
					AssetID:   0,
					Round:     txn.ConfirmedRound,
					Timestamp: time.Unix(int64(txn.RoundTime), 0),
				}, nil
			}
		}
	}

	// Check ASA payments if asset ID is not 0
	if payment.AssetID != 0 {
		asaResult, err := c.indexerClient.LookupAccountTransactions(payment.MerchantAddress).
			MinRound(lastCheckedRound + 1).
			TxType("axfer"). // For ASA transfers
			AssetID(payment.AssetID).
			Do(context.Background())

		if err != nil {
			return nil, fmt.Errorf("failed to lookup ASA transactions: %w", err)
		}

		for _, txn := range asaResult.Transactions {
			if txn.AssetTransferTransaction.Receiver == payment.MerchantAddress &&
				txn.AssetTransferTransaction.Amount >= payment.Amount {
				return &Transaction{
					ID:        txn.Id,
					Sender:    txn.Sender,
					Receiver:  txn.AssetTransferTransaction.Receiver,
					Amount:    txn.AssetTransferTransaction.Amount,
					AssetID:   payment.AssetID,
					Round:     txn.ConfirmedRound,
					Timestamp: time.Unix(int64(txn.RoundTime), 0),
				}, nil
			}
		}
	}

	return nil, nil // No matching transaction found
}

// GetLatestRound gets the latest round from the blockchain
func (c *Client) GetLatestRound() (uint64, error) {
	status, err := c.algodClient.Status().Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to get status: %w", err)
	}
	return status.LastRound, nil
}

// ValidateAddress validates an Algorand address
func (c *Client) ValidateAddress(address string) error {
	_, err := types.DecodeAddress(address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	return nil
}

// GetAssetInfo gets information about an ASA token
// func (c *Client) GetAssetInfo(assetID uint64) (types.AssetParams, error) {
// 	if assetID == 0 {
// 		// Return default params for ALGO
// 		return types.AssetParams{
// 			Total:     10000000000000000, // Total ALGO supply
// 			Decimals:  6,
// 			UnitName:  "ALGO",
// 			AssetName: "Algorand",
// 		}, nil
// 	}

// 	assetInfo, err := c.algodClient.GetAssetByID(assetID).Do(context.Background())
// 	if err != nil {
// 		return types.AssetParams{}, fmt.Errorf("failed to get asset info: %w", err)
// 	}
// 	return assetInfo.Params, nil
// }

// StartPaymentMonitor starts monitoring for payments
func (c *Client) StartPaymentMonitor(paymentChan chan<- *models.Payment, db PaymentDatabase) {
	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	defer ticker.Stop()

	var lastCheckedRound uint64 = 0

	for {
		select {
		case <-ticker.C:
			// Get current round
			currentRound, err := c.GetLatestRound()
			if err != nil {
				log.Printf("Error getting latest round: %v", err)
				continue
			}

			// Get pending payments
			payments, err := db.GetPendingPayments()
			if err != nil {
				log.Printf("Error getting pending payments: %v", err)
				continue
			}

			// Check each pending payment
			for _, payment := range payments {
				txn, err := c.CheckPayment(payment, lastCheckedRound)
				if err != nil {
					log.Printf("Error checking payment %s: %v", payment.ID, err)
					continue
				}

				if txn != nil {
					log.Printf("Payment found for %s: %s", payment.ID, txn.ID)
					payment.Status = models.PaymentStatusCompleted
					payment.TxnID = txn.ID
					payment.UpdatedAt = time.Now()

					// Send to payment channel for webhook processing
					paymentChan <- payment
				}
			}

			lastCheckedRound = currentRound
		}
	}
}

// PaymentDatabase interface for database operations
type PaymentDatabase interface {
	GetPendingPayments() ([]*models.Payment, error)
	UpdatePaymentStatus(id string, status models.PaymentStatus, txnID string) error
}
