package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// User is a DAO for users table
type User struct {
	ID      uuid.UUID       `db:"id"`
	Balance decimal.Decimal `db:"balance"`
}

// UserByID returns user by ID
func UserByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*User, error) {
	query := `SELECT id, balance FROM users WHERE id = $1`
	user := new(User)
	err := db.QueryRowxContext(ctx, query, id).StructScan(user)
	return user, err
}

// IncreaseUserBalance increases balance
func IncreaseUserBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, delta decimal.Decimal) error {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, delta, userID)
	return err
}

// DecreaseUserBalance decreases balance
func DecreaseUserBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, delta decimal.Decimal) error {
	query := `UPDATE users SET balance = balance - $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, delta, userID)
	return err
}

// ActualUserBalance returns user's balance
func ActualUserBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) (decimal.Decimal, error) {
	query := `SELECT balance FROM users WHERE id = $1`
	var balance decimal.Decimal
	err := tx.QueryRowxContext(ctx, query, userID).Scan(&balance)
	return balance, err
}
