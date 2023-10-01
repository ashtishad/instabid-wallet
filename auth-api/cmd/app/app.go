package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/ashtishad/instabid-wallet/auth-api/domain"
	"github.com/ashtishad/instabid-wallet/auth-api/service"
	"github.com/gin-gonic/gin"
)

func Start(srv *http.Server, dbClient *sql.DB, l *slog.Logger) {
	// Set GIN mode if specified in the environment variables
	if os.Getenv("GIN_MODE") != "" {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	// Create a new gin router
	var r = gin.New()
	srv.Handler = r

	// Wire up the handler for auth API
	authRepositoryDB := domain.NewAuthRepoDB(dbClient, l)
	ah := AuthHandlers{service.NewAuthService(authRepositoryDB, l)}

	// Route URL mappings for the auth API
	r.POST("/login", ah.LoginHandler)

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("could not start server: %v\n", "err", err.Error(), "srv", srv.Addr)
		}
	}()
}
