package middleware

import (
	"net/http"

	"github.com/pzabolotniy/user-balance/internal/logging"
)

// IsSourceTypeKnown validates "Source-Type" http header and declines request
// if passed value is not known
func IsSourceTypeKnown(knownSourceTypesList []string) func(next http.Handler) http.Handler {
	knownSourceTypes := make(map[string]bool, len(knownSourceTypesList))
	for i := range knownSourceTypesList {
		knownSourceTypes[knownSourceTypesList[i]] = true
	}
	fn := func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())
			headerSourceType := r.Header.Get("Source-Type")
			if _, ok := knownSourceTypes[headerSourceType]; !ok {
				logger.WithField("source_type", headerSourceType).Error("unknown source type")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(mw)
	}
	return fn
}
