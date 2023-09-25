package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ashtishad/instabid-wallet/lib"
	userAPI "github.com/ashtishad/instabid-wallet/user-api/cmd/app"
)

func main() {
	var wg sync.WaitGroup

	l := lib.InitSlogger()

	lib.SanityCheck(l)

	dbClient := lib.InitDB(l)
	defer func(dbClient *sql.DB) {
		if dbClsErr := dbClient.Close(); dbClsErr != nil {
			l.Error("unable to close db", "err", dbClsErr)
			os.Exit(1)
		}
	}(dbClient)

	userServer := lib.InitServerConfig("USER_API_PORT")

	wg.Add(1)

	go func() {
		userAPI.Start(userServer, dbClient, l)
		wg.Done()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	wg.Add(1)

	go lib.GracefulShutdown(ctx, userServer, &wg, "User")

	wg.Wait()
}
