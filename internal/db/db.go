package db

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pzabolotniy/user-balance/internal/config"
	"github.com/pzabolotniy/user-balance/internal/logging"
)

//NewDatabaseConnection set database configuration and initialize controller
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
	logger.
		WithField("addr",
			fmt.Sprintf("%v", conf.ConnString)).
		Debug("connected to database")

	return conn, err
}