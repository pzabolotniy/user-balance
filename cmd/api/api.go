package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pzabolotniy/user-balance/internal/api"
	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
	"github.com/pzabolotniy/user-balance/internal/migrations"
)

func main() {
	appConf, err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("init config failed: %s", err)
	}
	logger := logging.GetLogger()
	ctx := context.Background()
	ctx = logging.WithContext(ctx, logger)
	dbConn, err := db.Connect(ctx, appConf.Db)
	if err != nil {
		logger.WithError(err).Fatal("db connect failed, exiting")
	}

	err = migrations.Apply(ctx, appConf.Db)
	if err != nil {
		logger.WithError(err).Fatal("apply migrations failed, exiting")
	}

	router := chi.NewRouter()
	api.SetupRouter(router, logger, dbConn)
	if err := http.ListenAndServe(appConf.API.Bind, router); err != nil {
		logger.WithError(err).Error("api failed")
	}
}
