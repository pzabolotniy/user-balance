package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pzabolotniy/user-balance/internal/api"
	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
	"github.com/pzabolotniy/user-balance/internal/migrations"
	"github.com/pzabolotniy/user-balance/internal/txcancel"
)

func main() {
	appConf, err := config.GetAppConfig()
	if err != nil {
		log.Fatalf("init config failed: %s", err)
	}
	logger := logging.GetLogger()
	ctx := context.Background()
	ctx = logging.WithContext(ctx, logger)

	err = migrations.Apply(ctx, appConf.Db)
	if err != nil {
		logger.WithError(err).Fatal("apply migrations failed, exiting")
	}

	router := chi.NewRouter()
	dbConn, err := db.Connect(ctx, appConf.Db)
	if err != nil {
		logger.WithError(err).Error("db connect failed, exiting")
		return
	}
	api.SetupRouter(&api.SetupParams{
		Router:           router,
		Logger:           logger,
		DbConn:           dbConn,
		KnownSourceTypes: appConf.SourceTypes,
	})

	stop := make(chan bool)
	go func() {
		if err = http.ListenAndServe(appConf.API.Bind, router); err != nil {
			logger.WithError(err).Error("api listen failed")
		}
		stop <- true
	}()

	ticker := time.NewTicker(appConf.Cancelation.Interval)
	go func() {
		for { //nolint:gosimple
			select {
			case <-ticker.C:
				txcancel.CancelTransactions(ctx, appConf.Db, appConf.Cancelation)
			}
		}
	}()

	<-stop
}
