package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pzabolotniy/user-balance/internal/api"
	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

func main() {
	appConf := config.GetAppConfig()
	logger := logging.GetLogger()

	router := chi.NewRouter()
	api.SetupRouter(router, logger)
	if err := http.ListenAndServe(appConf.ServerConfig.Bind, router); err != nil {
		logger.WithError(err).Error("api failed")
	}
}
