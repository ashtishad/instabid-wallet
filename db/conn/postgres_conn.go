package conn

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// GetDsnURL constructs a PostgreSQL Data Source Name (DSN) URL using environment variables.
// It sets the connection parameters such as user, password, host, port, database name, timezone, and SSL mode.
// The resulting DSN URL is in the format:
// "postgres://user:password@host:port/dbname?sslmode=disable&timezone=UTC"
func GetDsnURL(l *slog.Logger) *url.URL {
	portInt, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		l.Error("error converting port string to int", "err", err.Error())
		os.Exit(1)
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWD")),
		Host:   fmt.Sprintf("%s:%d", os.Getenv("DB_HOST"), portInt),
		Path:   os.Getenv("DB_NAME"),
	}
	q := dsn.Query()
	q.Set("timezone", "UTC")
	q.Set("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	return &dsn
}

// GetDBClient creates a new database connection and returns it
func GetDBClient(l *slog.Logger) *sql.DB {
	connConfig, err := pgx.ParseConfig(GetDsnURL(l).String())

	if err != nil {
		l.Error("parsing postgres URI", "err", err)
		os.Exit(1)
	}

	db := stdlib.OpenDB(*connConfig)

	if err = db.Ping(); err != nil {
		l.Error("error pinging the database", "err", err.Error())
		os.Exit(1)
	}

	l.Info("successfully connected to database", "dsn", connConfig.ConnString())

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
