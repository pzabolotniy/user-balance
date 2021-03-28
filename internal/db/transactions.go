package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// Transaction is a DAO for transactions table
type Transaction struct {
	ID         uuid.UUID       `db:"id"`
	ExtID      string          `db:"external_tx_id"`
	UserID     uuid.UUID       `db:"user_id"`
	TxStateID  uuid.UUID       `db:"tx_state_id"`
	Amount     decimal.Decimal `db:"amount"`
	ReceivedAt time.Time       `db:"received_at"`
}

// TxByExternalID returns tx by external id
func TxByExternalID(ctx context.Context, db *sqlx.DB, txExtID string) (*Transaction, error) {
	query := `SELECT id, external_tx_id, user_id, tx_state_id, amount, received_at
FROM transactions
WHERE external_tx_id = $1`
	tx := new(Transaction)
	err := db.QueryRowxContext(ctx, query, txExtID).StructScan(tx)
	return tx, err
}

// CreateTransaction creates new transaction
func CreateTransaction(ctx context.Context, tx *sqlx.Tx, transaction *Transaction) error {
	query := `INSERT INTO transactions (
	id, external_tx_id, user_id, tx_state_id, amount, received_at
) VALUES (
	:id, :external_tx_id, :user_id, :tx_state_id, :amount, :received_at
)`
	_, err := tx.NamedExecContext(ctx, query, transaction)
	return err
}
