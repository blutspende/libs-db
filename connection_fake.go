package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type fakeDbConnection struct {
	debugLogEnabled bool
}

func (c *fakeDbConnection) SetSqlConnection(db *sqlx.DB) {
}

func (c *fakeDbConnection) EnableQueryLogging() {
	c.debugLogEnabled = true
}

func (c *fakeDbConnection) Ping() error {
	return nil
}

func (c *fakeDbConnection) BeginTx(ctx context.Context) (DbConnection, error) {
	return c, nil
}

func (c *fakeDbConnection) Commit() error {
	if c.debugLogEnabled {
		log.Debug().Msg("commit transaction")
	}
	return nil
}

func (c *fakeDbConnection) Rollback() error {
	if c.debugLogEnabled {
		log.Debug().Msg("rollback transaction")
	}
	return nil
}

func (c *fakeDbConnection) Rebind(query string) string {
	return query
}

func (c *fakeDbConnection) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return nil, nil
}

func (c *fakeDbConnection) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("ExecContext query")
	}
	return nil, nil
}

func (c *fakeDbConnection) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("GetContext query")
	}
	return nil
}

func (c *fakeDbConnection) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("arg", arg).Msg("NamedExecContext query")
	}
	return nil, nil
}

func (c *fakeDbConnection) NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", arg).Msg("NamedQueryContext query")
	}
	return nil, nil
}

func (c *fakeDbConnection) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("QueryxContext query")
	}
	return nil, nil
}

func (c *fakeDbConnection) QueryRowx(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("QueryRowxContext query")
	}
	return nil
}

// NewFakeDbConnection creates a new instance of a fake database connection that implements the DbConnection interface.
// All methods return nil value for sql.Result, *sqlx.Row, *sqlx.Rows, and nil error.
// Intended to be used in repository mocks in tests with service-layer transaction logic.
func NewFakeDbConnection() DbConnection {
	return &fakeDbConnection{}
}
