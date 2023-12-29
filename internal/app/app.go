package app

import (
	"context"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/internal/storage/psql"
	"github.com/migmatore/study-platform-api/internal/transport/rest"
	"github.com/migmatore/study-platform-api/internal/transport/rest/handler"
	"github.com/migmatore/study-platform-api/internal/usecase"
	"github.com/migmatore/study-platform-api/pkg/logger"
)

type App struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewApp(cfg *config.Config, logger logger.Logger) (App, error) {
	return App{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	a.logger.Info("Start app initializing...")

	a.logger.Info("Database connection initializing...")
	pool, err := psql.NewPostgres(ctx, 3, a.cfg, a.logger)
	if err != nil {
		a.logger.Fatalf("Failed to initialize db connection: %s", err.Error())
	}
	defer pool.Close()

	a.logger.Info("Database reconnection goroutine initializing...")
	go pool.Reconnect(ctx, a.cfg, a.logger)

	a.logger.Info("Storages initializing...")
	storages := storage.New(pool)

	a.logger.Info("Services initializing...")
	services := service.New(service.Deps{})

	a.logger.Info("Use cases initializing...")
	useCases := usecase.New(usecase.Deps{UserService: services.User})

	a.logger.Info("Handlers initializing...")
	restHandlers := handler.New(handler.Deps{AuthUseCase: useCases.Auth})

	app := restHandlers.Init(ctx)

	a.logger.Info("Server starting...")
	srv := rest.NewServer(":"+a.cfg.Server.Port, app, a.logger)
	srv.StartWithGracefulShutdown()
}
