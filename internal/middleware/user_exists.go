package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

type userIDctx struct{}

// UserExistsMw checks user existence
func UserExistsMw(dbConn *sqlx.DB) func(next http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			requestUserID := chi.URLParam(r, "user_id")
			userID, err := uuid.Parse(requestUserID)
			if err != nil {
				logger.
					WithError(err).
					WithField("request_user_id", requestUserID).
					Error("invalid user_id")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, err := db.UserByID(r.Context(), dbConn, userID)
			if err == nil {
				ctx := UserToContext(r.Context(), user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			logger.
				WithError(err).
				WithField("user_id", userID).
				Error("get user failed")
			w.WriteHeader(http.StatusNotFound)
		}
		return http.HandlerFunc(mw)
	}
	return fn
}

// UserToContext saves db user to context
func UserToContext(ctx context.Context, user *db.User) context.Context {
	return context.WithValue(ctx, userIDctx{}, user)
}

// UserFromContext extracts user from context
func UserFromContext(ctx context.Context) *db.User {
	logger := logging.FromContext(ctx)
	ctxUserID := ctx.Value(userIDctx{})
	if userID, ok := ctxUserID.(*db.User); ok {
		return userID
	}
	logger.
		WithFields(logging.Fields{
			"ctx_user_id": ctxUserID,
			"type":        fmt.Sprintf("%T", ctxUserID),
		}).
		Error("invalid user id type")
	var emptyUser db.User
	return &emptyUser
}
