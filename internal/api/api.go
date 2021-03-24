package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pzabolotniy/user-balance/internal/logging"
	"github.com/pzabolotniy/user-balance/internal/middleware"
)

// SetupRouter setup passed gin-router (*gin.Engine)
// to prepare http server
func SetupRouter(router *chi.Mux, logger logging.Logger) {
	env := NewEnv()

	router.Use(middleware.WithLoggerMw(logger))
	router.Use(middleware.WithUniqRequestID)
	router.Use(middleware.LogRequestBoundariesMw)

	router.Route("/api/v1", func(versionedRoter chi.Router) {
		versionedRoter.Route("/users/{user_id}", func(userRoute chi.Router) {
			userRoute.Route("/transactions/{tx_id}", func(txRouter chi.Router) {
				txRouter.Post("/", env.PostUserTransaction)
			})
		})
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
