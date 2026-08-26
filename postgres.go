package db

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

const MaxQueryParams = 65534

type PgConfig struct {
	ApplicationName              string
	Host                         string
	Port                         uint32
	User                         string
	Pass                         string
	Database                     string
	SSLMode                      string
	MaxOpenConnections           *int
	MaxIdleConnections           *int
	ConnectionMaxLifetimeSeconds *int
	ConnectionMaxIdleTimeSeconds *int
	UseOpenTelemetry             bool
}

type Postgres interface {
	Connect(ctx context.Context) (*sqlx.DB, error)
	GetSqlConnection() (*sqlx.DB, error)
	Close() error
}

type postgres struct {
	config PgConfig
	pgConn *sqlx.DB
}

func NewPostgres(config PgConfig) Postgres {
	return &postgres{
		config: config,
		pgConn: nil,
	}
}

func (p *postgres) Connect(ctx context.Context) (pgDB *sqlx.DB, err error) {
	// Skip if already connected
	if p.pgConn != nil {
		err = p.pgConn.Ping()
		if err == nil {
			log.Warn().Msg("postgres already connected")
			return p.pgConn, nil
		}
	}
	// Setup connection parameters
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s application_name=%s",
		p.config.Host, p.config.Port, p.config.User,
		p.config.Pass, p.config.Database, p.config.SSLMode, p.config.ApplicationName)
	driverName := "pgx"
	// Setup OpenTelemetry if enabled
	if p.config.UseOpenTelemetry {
		driverName, err = otelsql.Register(driverName,
			otelsql.AllowRoot(),
			otelsql.TraceQueryWithoutArgs(),
			otelsql.TraceRowsClose(),
			otelsql.TraceRowsAffected(),
			otelsql.WithDatabaseName(p.config.Database),
			otelsql.WithSystem(semconv.DBSystemNamePostgreSQL),
		)
		if err != nil {
			log.Error().Err(err).Msg("open telemetry setup failed")
			return nil, err
		}
		// Note: must be explicitly specified, because otelsql changes the driver name, which breaks sqlx parameter binding in named queries
		sqlx.BindDriver(driverName, sqlx.DOLLAR)
	}
	// Connect to Postgres
	pgDB, err = sqlx.ConnectContext(ctx, driverName, url)
	if err != nil {
		log.Error().Err(err).Msg("connecting to postgres failed")
		return nil, err
	}
	if pgDB == nil {
		log.Error().Msg("connecting to postgres failed")
		return nil, err
	}
	p.pgConn = pgDB
	// Configure connection pool
	if p.config.MaxOpenConnections != nil {
		pgDB.DB.SetMaxOpenConns(*p.config.MaxOpenConnections)
	}
	if p.config.MaxIdleConnections != nil {
		pgDB.DB.SetMaxIdleConns(*p.config.MaxIdleConnections)
	}
	if p.config.ConnectionMaxLifetimeSeconds != nil {
		pgDB.DB.SetConnMaxLifetime(time.Duration(*p.config.ConnectionMaxLifetimeSeconds) * time.Second)
	}
	if p.config.ConnectionMaxIdleTimeSeconds != nil {
		pgDB.DB.SetConnMaxIdleTime(time.Duration(*p.config.ConnectionMaxIdleTimeSeconds) * time.Second)
	}
	// Log successful connection and return
	log.Info().Msgf("postgres available, connected to %s / %s", p.config.Host, p.config.Database)
	return p.pgConn, nil
}

func (p *postgres) GetSqlConnection() (*sqlx.DB, error) {
	if p.pgConn == nil {
		return nil, ErrNoPgConnection
	}
	return p.pgConn, nil
}

func (p *postgres) Close() error {
	if p.pgConn != nil {
		err := p.pgConn.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close postgres connection")
			return err
		}
		p.pgConn = nil
		log.Info().Msg("postgres connection closed")
	}
	return nil
}
