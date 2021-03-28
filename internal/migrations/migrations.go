// Package migrations - database migrations utilities
package migrations

import (
	"context"

	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/db"
	"github.com/pzabolotniy/user-balance/internal/logging"
	migrate "github.com/rubenv/sql-migrate"
)

// Apply applies database migrations
func Apply(ctx context.Context, conf *config.DB) error {
	logger := logging.FromContext(ctx)

	logger.Debug("Applying migrations")

	dbConn, err := db.Connect(ctx, conf)
	if err != nil {
		logger.WithError(err).Fatal("db connect failed")
	}
	defer func() {
		clErr := dbConn.Close()
		if clErr != nil {
			logger.WithError(clErr).Error("closing db connection is failed")
		}
	}()

	migrationDirection := migrate.Up
	migrationCount := -1
	migrate.SetTable(conf.MigrationTable)
	count, err := migrate.ExecMax(dbConn.DB, "postgres",
		&migrate.FileMigrationSource{Dir: conf.MigrationDirPath},
		migrationDirection, migrationCount,
	)
	if err != nil {
		logger.Fatal(err)
	}
	logger.WithField("count", count).Info("migrations applied")
	return err
}
