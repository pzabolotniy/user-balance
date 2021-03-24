package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

// WithLoggerMw injects logger to the request context
func WithLoggerMw(logger logging.Logger) func(next http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx := logging.WithContext(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(mw)
	}
	return fn
}

// LogRequestBoundariesMw logs start and end of the request
func LogRequestBoundariesMw(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		uri := r.URL.String()
		logger.WithField("path", uri).Trace("REQUEST STARTED")
		next.ServeHTTP(w, r)
		logger.Trace("REQUEST FINISHED")
	}
	return http.HandlerFunc(mw)
}

// WithUniqRequestID appends uuid to the logger for every request
func WithUniqRequestID(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		uniqRequestUUID := uuid.New()
		w.Header().Set("X-Request-Id", uniqRequestUUID.String())
		logger = logger.WithField("x-request-id", uniqRequestUUID.String())
		ctx := logging.WithContext(r.Context(), logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(mw)
}
