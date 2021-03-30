package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// CanceledTx is a DAO for canceled_txs row
type CanceledTx struct {
	ID         uuid.UUID `db:"id"`
	TxID       uuid.UUID `db:"tx_id"`
	CanceledAt time.Time `db:"canceled_at"`
}

// CreateCanceledTx creates new row into canceled_txs table
func CreateCanceledTx(ctx context.Context, tx *sqlx.Tx, canceledTx *CanceledTx) error {
	query := `INSERT INTO canceled_txs (
	id, tx_id, canceled_at
) VALUES (
	:id, :tx_id, :canceled_at
)`
	_, err := tx.NamedExecContext(ctx, query, canceledTx)
	return err
}
