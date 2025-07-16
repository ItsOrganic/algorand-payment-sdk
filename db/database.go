package db

import (
	"database/sql"
	"fmt"

	"algopay/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return database, nil
}

// createTables creates the necessary database tables
func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS payments (
		id TEXT PRIMARY KEY,
		merchant_address TEXT NOT NULL,
		amount INTEGER NOT NULL,
		asset_id INTEGER NOT NULL DEFAULT 0,
		callback_url TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		txn_id TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
	CREATE INDEX IF NOT EXISTS idx_payments_merchant ON payments(merchant_address);
	CREATE INDEX IF NOT EXISTS idx_payments_expires ON payments(expires_at);
	`
	_, err := d.db.Exec(query)
	return err
}

// CreatePayment creates a new payment record
func (d *Database) CreatePayment(payment *models.Payment) error {
	query := `
	INSERT INTO payments (id, merchant_address, amount, asset_id, callback_url, status, created_at, updated_at, expires_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query,
		payment.ID,
		payment.MerchantAddress,
		payment.Amount,
		payment.AssetID,
		payment.CallbackURL,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.ExpiresAt,
	)
	return err
}

// GetPayment retrieves a payment by ID
func (d *Database) GetPayment(id string) (*models.Payment, error) {
	query := `
	SELECT id, merchant_address, amount, asset_id, callback_url, status, txn_id, created_at, updated_at, expires_at
	FROM payments WHERE id = ?
	`
	row := d.db.QueryRow(query, id)

	payment := &models.Payment{}
	var txnID sql.NullString
	err := row.Scan(
		&payment.ID,
		&payment.MerchantAddress,
		&payment.Amount,
		&payment.AssetID,
		&payment.CallbackURL,
		&payment.Status,
		&txnID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	if txnID.Valid {
		payment.TxnID = txnID.String
	}

	return payment, nil
}

// UpdatePaymentStatus updates the payment status and transaction ID
func (d *Database) UpdatePaymentStatus(id string, status models.PaymentStatus, txnID string) error {
	query := `
	UPDATE payments 
	SET status = ?, txn_id = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := d.db.Exec(query, status, txnID, id)
	return err
}

// GetPendingPayments retrieves all pending payments
func (d *Database) GetPendingPayments() ([]*models.Payment, error) {
	query := `
	SELECT id, merchant_address, amount, asset_id, callback_url, status, txn_id, created_at, updated_at, expires_at
	FROM payments 
	WHERE status = 'pending' AND expires_at > CURRENT_TIMESTAMP
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		payment := &models.Payment{}
		var txnID sql.NullString
		err := rows.Scan(
			&payment.ID,
			&payment.MerchantAddress,
			&payment.Amount,
			&payment.AssetID,
			&payment.CallbackURL,
			&payment.Status,
			&txnID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}

		if txnID.Valid {
			payment.TxnID = txnID.String
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

// ExpireOldPayments marks expired payments as expired
func (d *Database) ExpireOldPayments() error {
	query := `
	UPDATE payments 
	SET status = 'expired', updated_at = CURRENT_TIMESTAMP
	WHERE status = 'pending' AND expires_at <= CURRENT_TIMESTAMP
	`
	_, err := d.db.Exec(query)
	return err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
