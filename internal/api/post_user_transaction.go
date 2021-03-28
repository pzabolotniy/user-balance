package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
	"github.com/pzabolotniy/user-balance/internal/middleware"
	"github.com/shopspring/decimal"
)

type TxPayload struct {
	State         string `json:"state" validate:"required,oneof=win lost"`
	Amount        string `json:"amount" validate:"required"`
	TransactionID string `json:"transactionId" validate:"required"`
}

// PostUserTransaction handles users' transactions
func (env *Env) PostUserTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)

	validate := validator.New()
	txPayload := new(TxPayload)
	if err := json.NewDecoder(r.Body).Decode(txPayload); err != nil {
		logger.WithError(err).Error("read input failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := validate.Struct(txPayload); err != nil {
		logger.WithError(err).Error("validate input failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	txAmount, err := decimal.NewFromString(txPayload.Amount)
	if err != nil {
		logger.WithError(err).Error("bad amount")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	txState := ""
	switch txPayload.State {
	case "win":
		txState = db.Win
	case "lost":
		txState = db.Lost
	}
	if txState == "" {
		logger.WithField("state", txPayload.State).Error("unknown state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbUser := middleware.UserFromContext(ctx)
	dbConn := env.DbConn
	tx, err := dbConn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		logger.WithError(err).Error("start transaction failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		txErr := tx.Rollback()
		if txErr != nil && txErr != sql.ErrTxDone {
			logger.WithError(err).Error("rollback transaction failed")
		}
	}()

	// we can not rely on the dbUser.Amount value, because it could be changed
	// before the transaction started
	actualBalance, err := db.ActualUserBalance(ctx, tx, dbUser.ID)
	if err != nil {
		logger.
			WithError(err).
			WithField("user_id", dbUser.ID).
			Error("get actual balance failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbTxState, err := db.TxStateByName(ctx, tx, txState)
	if err != nil {
		logger.
			WithError(err).
			WithField("state", txState).
			Error("get tx state failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	txID := uuid.New()
	receivedAt := time.Now().UTC()
	txExtID := chi.URLParam(r, "tx_id")
	newDbTransaction := &db.Transaction{
		ID:         txID,
		ExtID:      txExtID,
		UserID:     dbUser.ID,
		TxStateID:  dbTxState.ID,
		Amount:     txAmount,
		ReceivedAt: receivedAt,
	}
	err = db.CreateTransaction(ctx, tx, newDbTransaction)
	if err != nil {
		logger.WithError(err).Error("create transaction failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch txState {
	case db.Win:
		err = db.IncreaseUserBalance(ctx, tx, dbUser.ID, txAmount)
	case db.Lost:
		if actualBalance.Sub(txAmount).IsNegative() {
			logger.Error("balance can not be negative, can not decrease it")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.DecreaseUserBalance(ctx, tx, dbUser.ID, txAmount)
	default:
		err = fmt.Errorf("unknown tx state '%s'", txState)
	}

	if err != nil {
		logger.WithError(err).Error("update balance failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("commit failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
