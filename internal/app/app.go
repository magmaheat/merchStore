package app

import (
	"fmt"
	"github.com/magmaheat/merchStore/internal/config"
	v1 "github.com/magmaheat/merchStore/internal/http/v1"
	"github.com/magmaheat/merchStore/internal/repo/pgdb"
	"github.com/magmaheat/merchStore/internal/service"
	"github.com/magmaheat/merchStore/pkg/httpserver"
	"github.com/magmaheat/merchStore/pkg/postgres"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	cfg := config.New(configPath)

	setupLogger(cfg.Log.Level)

	log.Info("Initializing postgres")
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pg.Close()

	log.Info("Initializing storage...")
	storage := pgdb.New(pg)

	log.Info("Initializing services")
	deps := service.Dependencies{
		Repo:     storage,
		SignKey:  cfg.SigningKey,
		TokenTTL: cfg.TokenAccessTTL,
	}
	services := service.NewService(deps)

	log.Info("Initializing handlers and routes...")
	routes := v1.New(services)

	log.Info("Starting http server")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(routes, httpserver.Port(cfg.HTTP.Port))

	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
