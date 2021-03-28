package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	Win  = "win"
	Lost = "lost"
)

type TxState struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func TxStateByName(ctx context.Context, txConn *sqlx.Tx, name string) (*TxState, error) {
	query := `SELECT id, name FROM tx_states WHERE name = $1`
	txState := new(TxState)
	err := txConn.QueryRowxContext(ctx, query, name).StructScan(txState)
	return txState, err
}
