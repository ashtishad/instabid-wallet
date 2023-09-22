package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ashtishad/ecommerce/db/conn"
	"github.com/golang-migrate/migrate/v4"
)

func InitServerConfig(portEnv string) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv(portEnv)),
		IdleTimeout:    100 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func InitSlogger() *slog.Logger {
	handlerOpts := GetSlogConf()
	l := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(l)

	return l
}

// InitDB initializes db connections, applies migrations and execute any bulk insert functions.
func InitDB(l *slog.Logger) *sql.DB {
	dbClient := conn.GetDBClient(l)

	m, err := migrate.New(
		"file://db/migrations",
		conn.GetDSNString(l).String(),
	)
	if err != nil {
		l.Error("error creating migration: %v", "err", err.Error())
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error("error applying migration: %v", "err", err.Error())
	}

	return dbClient
}

// SanityCheck checks that all required environment variables are set.
// If any of the required variables is not defined, it sets a default value and prints a warning message.
// If error happened during setting env variable, then logs error and exits application.
func SanityCheck(l *slog.Logger) {
	defaultEnvVars := map[string]string{
		"API_HOST":  "127.0.0.1",
		"API_PORT":  "8000",
		"DB_USER":   "postgres",
		"DB_PASSWD": "postgres",
		"DB_HOST":   "127.0.0.1",
		"DB_PORT":   "5432",
		"DB_NAME":   "instabid",
	}

	for key, defaultValue := range defaultEnvVars {
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, defaultValue); err != nil {
				l.Error(fmt.Sprintf("failed to set environment variable %s to default value %s. Exiting application.", key, defaultValue))
				os.Exit(1)
			}

			l.Warn(fmt.Sprintf("environment variable %s not defined. Setting to default: %s", key, defaultValue))
		}
	}
}
