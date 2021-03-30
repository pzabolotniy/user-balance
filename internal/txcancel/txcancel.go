package txcancel

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

// CancelTransactions contains main transactions' cancelation
//nolint:funlen
func CancelTransactions(ctx context.Context, dbConnConf *config.DB, cancelationConf *config.Cancelation) {
	logger := logging.FromContext(ctx)
	dbConn, err := db.Connect(ctx, dbConnConf)
	if err != nil {
		logger.WithError(err).Error("db connect failed, exiting")
		return
	}
	defer func() {
		if err = db.Disconnect(ctx, dbConn); err != nil {
			logger.WithError(err).Error("disconnect failed")
		}
	}()

	tx, err := dbConn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		logger.WithError(err).Error("start transaction failed, exiting")
		return
	}
	defer func() {
		txErr := tx.Rollback()
		if txErr != nil && txErr != sql.ErrTxDone {
			logger.WithError(err).Error("rollback transaction failed")
		}
	}()

	statesList, err := db.GetStates(ctx, tx)
	if err != nil {
		logger.WithError(err).Error("get states failed")
		return
	}
	states := stateListToMap(statesList)

	// multiply count, because later will take only odd txs
	transactionsToCancel, err := db.GetTxsToCancel(ctx, tx, cancelationConf.TxsPerRound*2) //nolint:gomnd
	if err != nil {
		logger.WithError(err).Error("get transactions to cancel failed")
		return
	}
	oddTxs := getOddTxs(transactionsToCancel)
	logger.WithField("txs_to_cancel", len(oddTxs)).Info("transactions will be canceled")

	for i := range oddTxs {
		txToCancel := oddTxs[i]
		txLogger := logger.WithField("tx_id", txToCancel.ID)

		txLogger.Trace("canceling transaction")
		switch states[txToCancel.TxStateID] {
		case db.Win: // inverting tx
			err = db.DecreaseUserBalance(ctx, tx, txToCancel.UserID, txToCancel.Amount)
		case db.Lost: // inverting tx
			err = db.IncreaseUserBalance(ctx, tx, txToCancel.UserID, txToCancel.Amount)
		default:
			err = fmt.Errorf("unknown tx state id '%s'", txToCancel.TxStateID)
		}
		if err != nil {
			txLogger.
				WithError(err).
				Error("cancel tx failed")
			return
		}

		canceledTxID := uuid.New()
		canceledTx := &db.CanceledTx{
			ID:         canceledTxID,
			TxID:       txToCancel.ID,
			CanceledAt: time.Now().UTC(),
		}
		err = db.CreateCanceledTx(ctx, tx, canceledTx)
		if err != nil {
			txLogger.WithError(err).Error("save canceled failed")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("commit failed")
	}
}

func stateListToMap(statesList []db.TxState) map[uuid.UUID]string {
	states := make(map[uuid.UUID]string, len(statesList))
	for _, state := range statesList {
		states[state.ID] = state.Name
	}
	return states
}

func getOddTxs(list []db.TransactionToCancel) []db.TransactionToCancel {
	oddList := make([]db.TransactionToCancel, 0)
	for i := range list {
		if i%2 == 1 {
			if list[i].CanceledID == nil { // transaction is not canceled yet
				oddList = append(oddList, list[i])
			}
		}
	}
	return oddList
}
