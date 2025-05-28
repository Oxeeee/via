package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/OxytocinGroup/theca-v3/internal/config"
	"github.com/OxytocinGroup/theca-v3/internal/database"
	"github.com/OxytocinGroup/theca-v3/internal/model"
	"github.com/OxytocinGroup/theca-v3/internal/repository"
	"github.com/OxytocinGroup/theca-v3/internal/server"
	"github.com/OxytocinGroup/theca-v3/internal/server/handlers"
	"github.com/OxytocinGroup/theca-v3/internal/server/middleware"
	"github.com/OxytocinGroup/theca-v3/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/OxytocinGroup/theca-v3/docs"
)

type Application struct {
	cfg            *config.Config
	log            *slog.Logger
	server         *server.Server
	authMiddleware middleware.AuthMiddleware
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger) *Application {
	server := server.New(cfg, log)

	if cfg.IsLocalRun {
		server.Router().Use(gin.Logger())
	}

	server.Router().Use(middleware.MetricsMiddleware())

	db, err := database.ConnectDatabase(ctx, cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Bookmark{}); err != nil {
		log.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	repo := repository.NewRepository(db.GetDB(), log)
	service := service.NewService(repo, log, cfg)

	handlers := handlers.NewHandler(service, log)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)

	initHandlers(server, handlers, authMiddleware)
	initPrivateHandlers(server)

	app := &Application{
		cfg:            cfg,
		log:            log,
		server:         server,
		authMiddleware: authMiddleware,
	}

	return app
}

func initHandlers(server *server.Server, handlers *handlers.Handler, authMiddleware middleware.AuthMiddleware) {
	v1 := server.Router().Group("/v1")
	v1.POST("/register", handlers.Register)
	v1.POST("/login", handlers.Login)

	sec := v1.Group("/api", authMiddleware.JWTMiddleware())
	sec.DELETE("/logout", handlers.Logout)

}

func initPrivateHandlers(server *server.Server) {
	server.PrivateRouter().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	server.PrivateRouter().GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func (a *Application) Run() {
	const op = "app.Run"
	a.server.Start()
	log := a.log.With(slog.String("op", op))
	log.Info("application started",
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
		slog.String("name", a.cfg.AppName),
	)
}

func (a *Application) Stop() {
	const op = "app.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("shutting down application...")

	a.server.Stop()

	a.log.Info("application stopped")
}
