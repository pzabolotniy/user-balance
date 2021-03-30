package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pzabolotniy/user-balance/internal/logging"
	"github.com/pzabolotniy/user-balance/internal/middleware"
)

// SetupParams contains settings to setup api router
type SetupParams struct {
	Router           *chi.Mux
	Logger           logging.Logger
	DbConn           *sqlx.DB
	KnownSourceTypes []string
}

// SetupRouter setup passed gin-router (*gin.Engine)
// to prepare http server
func SetupRouter(params *SetupParams) {
	env := NewEnv(WithDbConn(params.DbConn))
	router := params.Router

	router.Use(middleware.WithLoggerMw(params.Logger))
	router.Use(middleware.WithUniqRequestID)
	router.Use(middleware.LogRequestBoundariesMw)
	router.Use(middleware.IsSourceTypeKnown(params.KnownSourceTypes))

	router.Route("/api/v1", func(versionedRoter chi.Router) {
		versionedRoter.Route("/users/{user_id}", func(userRoute chi.Router) {
			userRoute.Use(middleware.UserExistsMw(params.DbConn))
			userRoute.Route("/transactions/{tx_id}", func(txRouter chi.Router) {
				txRouter.Use(middleware.TransactionExistsMw(params.DbConn))
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
