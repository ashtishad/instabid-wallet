package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/ashtishad/instabid-wallet/user-api/internal/service"
	"github.com/gin-gonic/gin"
)

func Start(srv *http.Server, dbClient *sql.DB, l *slog.Logger) {
	gin.SetMode(gin.ReleaseMode)

	var r = gin.New()
	srv.Handler = r

	// wire up the handler
	userRepositoryDB := domain.NewUserRepoDB(dbClient, l)
	uh := UserHandlers{service.NewUserService(userRepositoryDB)}

	// route url mappings
	setUsersAPIRoutes(r, uh)

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("could not start server: %v\n", "err", err.Error(), "srv", srv.Addr)
		}
	}()
}

func setUsersAPIRoutes(r *gin.Engine, uh UserHandlers) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", uh.NewUserHandler)
	}
}
