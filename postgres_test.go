package db

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresConnection(t *testing.T) {
	// General setup
	var err error
	ctx := context.Background()

	// Setup test container
	dbName := "test"
	dbUser := "user"
	dbPass := "pass"
	postgresContainer, err := tcpg.Run(ctx,
		"postgres:16-alpine",
		tcpg.WithDatabase(dbName),
		tcpg.WithUsername(dbUser),
		tcpg.WithPassword(dbPass),
		tcpg.BasicWaitStrategies(),
	)
	assert.Nil(t, err)
	defer func() {
		err = testcontainers.TerminateContainer(postgresContainer)
		assert.Nil(t, err)
	}()

	// Setup postgres instance
	connStr, err := postgresContainer.ConnectionString(ctx)
	assert.Nil(t, err)
	log.Info().Str("connStr", connStr).Send()
	dbPort, err := postgresContainer.MappedPort(ctx, "5432")
	assert.Nil(t, err)
	dbConfig := PgConfig{
		ApplicationName: "test",
		Host:            "localhost",
		Port:            (uint32)(dbPort.Num()),
		User:            dbUser,
		Pass:            dbPass,
		Database:        dbName,
		SSLMode:         "disable",
	}
	pg := NewPostgres(dbConfig)

	// Test connection before Connect() is called
	dbConn, err := pg.GetSqlConnection()
	assert.Nil(t, dbConn)
	assert.ErrorIs(t, err, ErrNoPgConnection)

	// Test Connect() method
	dbConn, err = pg.Connect(context.Background())
	assert.NotNil(t, dbConn)
	assert.Nil(t, err)

	// Test GetSqlConnection() after Connect() is called
	dbConn, err = pg.GetSqlConnection()
	assert.NotNil(t, dbConn)
	assert.Nil(t, err)

	// Test Connect() method when already connected
	dbConn, err = pg.Connect(context.Background())
	assert.NotNil(t, dbConn)
	assert.Nil(t, err)

	// Test Close() method
	err = pg.Close()
	assert.Nil(t, err)

	// Test GetSqlConnection() after Close() is called
	dbConn, err = pg.GetSqlConnection()
	assert.Nil(t, dbConn)
	assert.ErrorIs(t, err, ErrNoPgConnection)
}
