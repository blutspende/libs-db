# db
Contains the `Postgres` class to handle Postgres connection and some utility functions for SQL specific nullable types.

###### Install
`go get github.com/blutspende/libs-db`

## Postgres
`Postgres` is a class for handling Postgres connections. It provides methods for connecting, disconnecting, and obtaining the underlying raw SQL connection `*sqlx.DB`.
`NewPostgres` is used to create a new instance, which requires a `PgConfig` as input:
```go
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
}
```
The `*int` types can be set to `nil` to avoid setting those configurations on the database connection.

## DbConnection
`DbConnection` is a class for transaction and query handling. It allows direct execution of queries, as well as transaction management with `BeginTx`, `Commit` and `Rollback` methods.
Specific error codes and code conversions are also provided.

## Utility functions
```go
func NullStringToString(value sql.NullString) string
func NullStringToStringPointer(value sql.NullString) *string
func NullTimeToTimePointer(value sql.NullTime) *time.Time
func TimePointerToNullTime(value *time.Time) sql.NullTime
```