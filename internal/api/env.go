package api

import "github.com/jmoiron/sqlx"

// Env is a container for api
// environment variables
type Env struct {
	DbConn *sqlx.DB
}

// OptionFunc is a type of args for the NewEnv
// this funcs are called in the constructor
// to init Env struct
type OptionFunc func(e *Env)

// NewEnv is a constructor for the *Env
// *Env has no default values
func NewEnv(opts ...OptionFunc) *Env {
	env := new(Env)
	for _, optFunc := range opts {
		optFunc(env)
	}

	return env
}

func WithDbConn(db *sqlx.DB) OptionFunc {
	return func(e *Env) {
		e.DbConn = db
	}
}
