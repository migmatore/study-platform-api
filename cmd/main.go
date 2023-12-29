package main

import (
	"context"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/internal/app"
	"github.com/migmatore/study-platform-api/pkg/logger"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting api server")

	cfgFile, err := config.LoadConfig("./config/config-local")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("LogLevel: %s, Mode: %s", cfg.Logger.Level, cfg.Server.Mode)

	a, err := app.NewApp(cfg, appLogger)
	if err != nil {
		appLogger.Fatal(err)
	}

	a.Run(ctx)
}
