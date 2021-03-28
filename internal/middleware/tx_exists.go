package middleware

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

func TransactionExistsMw(dbConn *sqlx.DB) func(next http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			requestTxID := chi.URLParam(r, "tx_id")
			_, err := db.TxByExternalID(r.Context(), dbConn, requestTxID)
			if err != sql.ErrNoRows {
				if err == nil {
					logger.
						WithField("tx_id", requestTxID).
						Error("transaction processed")
					w.WriteHeader(http.StatusConflict)
					return
				}
				logger.
					WithError(err).
					WithField("tx_id", requestTxID).
					Error("get transaction failed")
				w.WriteHeader(http.StatusInternalServerError)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
	return fn
}
