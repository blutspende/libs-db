package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type DbConnection interface {
	SetSqlConnection(db *sqlx.DB)
	EnableQueryLogging()
	Ping() error
	BeginTx(ctx context.Context) (DbConnection, error)
	Commit() error
	Rollback() error
	Rebind(query string) string
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type dbConnection struct {
	db              *sqlx.DB
	tx              *sqlx.Tx
	debugLogEnabled bool
}

func NewEmptyDbConnection() DbConnection {
	return &dbConnection{}
}
func NewDbConnection(db *sqlx.DB) DbConnection {
	return &dbConnection{
		db:              db,
		tx:              nil,
		debugLogEnabled: false,
	}
}

func (c *dbConnection) SetSqlConnection(db *sqlx.DB) {
	c.db = db
}
func (c *dbConnection) EnableQueryLogging() {
	c.debugLogEnabled = true
}

func (c *dbConnection) Ping() error {
	return c.db.Ping()
}

func (c *dbConnection) BeginTx(ctx context.Context) (DbConnection, error) {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg(ErrBeginTransactionFailed.Error())
		return nil, ErrBeginTransactionFailed
	}
	connCopy := *c
	connCopy.tx = tx
	return &connCopy, err
}

func (c *dbConnection) Commit() error {
	if c.debugLogEnabled {
		log.Debug().Msg("commit transaction")
	}
	if c.tx == nil {
		return ErrCommitWithoutTransaction
	}
	err := c.tx.Commit()
	c.tx = nil
	if err != nil {
		log.Error().Err(err).Msg(ErrCommitTransactionFailed.Error())
		return ErrCommitTransactionFailed
	}
	return nil
}

func (c *dbConnection) Rollback() error {
	if c.debugLogEnabled {
		log.Debug().Msg("rollback transaction")
	}
	if c.tx == nil {
		return ErrRollbackWithoutTransaction
	}
	err := c.tx.Rollback()
	c.tx = nil
	if err != nil {
		log.Error().Err(err).Msg(ErrRollbackTransactionFailed.Error())
		return ErrRollbackTransactionFailed
	}
	return nil
}

func (c *dbConnection) Rebind(query string) string {
	if c.tx != nil {
		return c.tx.Rebind(query)
	}
	return c.db.Rebind(query)
}

func (c *dbConnection) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	if c.tx != nil {
		return c.tx.PrepareNamed(query)
	}
	return c.db.PrepareNamed(query)
}

func (c *dbConnection) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("ExecContext query")
	}
	if c.tx != nil {
		return c.tx.ExecContext(ctx, query, args...)
	}
	return c.db.ExecContext(ctx, query, args...)
}

func (c *dbConnection) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("GetContext query")
	}
	if c.tx != nil {
		return c.tx.GetContext(ctx, dest, query, args...)
	}
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c *dbConnection) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("arg", arg).Msg("NamedExecContext query")
	}
	if c.tx != nil {
		return c.tx.NamedExecContext(ctx, query, arg)
	}
	return c.db.NamedExecContext(ctx, query, arg)
}

func (c *dbConnection) NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", arg).Msg("NamedQueryContext query")
	}
	if c.tx != nil {
		return sqlx.NamedQueryContext(ctx, c.tx, query, arg)
	}
	return c.db.NamedQueryContext(ctx, query, arg)
}

func (c *dbConnection) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("QueryxContext query")
	}
	if c.tx != nil {
		return c.tx.QueryxContext(ctx, query, args...)
	}
	return c.db.QueryxContext(ctx, query, args...)
}

func (c *dbConnection) QueryRowx(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	if c.debugLogEnabled {
		log.Debug().Ctx(ctx).Str("query", query).Interface("args", args).Msg("QueryRowxContext query")
	}
	if c.tx != nil {
		return c.tx.QueryRowxContext(ctx, query, args...)
	}
	return c.db.QueryRowxContext(ctx, query, args...)
}
