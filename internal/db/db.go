package db

import (
	"context"

	// import postgresql driver
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

// Connect creates new db connection
func Connect(ctx context.Context, conf *config.DB) (*sqlx.DB, error) {
	var logger = logging.FromContext(ctx)
	logger.WithField("conn_string", conf.ConnString).Trace("connecting to db")
	var conn, err = sqlx.ConnectContext(ctx, "pgx", conf.ConnString)
	if err != nil {
		logger.WithError(err).Error("unable to connect to database")
		return nil, err
	}
	conn.DB.SetMaxOpenConns(conf.MaxOpenConns)
	conn.DB.SetConnMaxLifetime(conf.ConnMaxLifetime)

	return conn, err
}
