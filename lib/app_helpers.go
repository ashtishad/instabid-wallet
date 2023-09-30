package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ashtishad/instabid-wallet/db/conn"
	"github.com/golang-migrate/migrate/v4"

	// ignore: revive
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		conn.GetDsnURL(l).String(),
	)

	if err != nil {
		l.Error("error creating migration", "err", err.Error())
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error("error applying migration", "err", err.Error())
	}

	return dbClient
}

// SanityCheck checks that all required environment variables are set.
// If any of the required variables is not defined, it sets a default value and prints a warning message.
// If error happened during setting env variable, then logs error and exits application.
func SanityCheck(l *slog.Logger) {
	defaultEnvVars := map[string]string{
		"API_HOST":      "127.0.0.1",
		"USER_API_PORT": "8000",
		"AUTH_API_PORT": "8001",
		"DB_USER":       "postgres",
		"DB_PASSWD":     "postgres",
		"DB_HOST":       "127.0.0.1",
		"DB_PORT":       "5432",
		"DB_NAME":       "instabid",
		"GIN_MODE":      "debug",
		"HMACSecret":    "hmacSampleSecret",
	}

	for key, defaultValue := range defaultEnvVars {
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, defaultValue); err != nil {
				l.Error(fmt.Sprintf(
					"failed to set environment variable %s to default value %s. Exiting application.",
					key,
					defaultValue,
				))
				os.Exit(1)
			}

			l.Warn(fmt.Sprintf("environment variable %s not defined. Setting to default: %s", key, defaultValue))
		}
	}
}

func GracefulShutdown(ctx context.Context, srv *http.Server, wg *sync.WaitGroup, serverName string) {
	defer wg.Done()
	log.Printf("Shutting down %s server...\n", serverName)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully shutdown the %s server: %v\n", serverName, err)
	}

	log.Printf("%s server gracefully stopped\n", serverName)
}
