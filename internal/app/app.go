package app

import (
	"github.com/magmaheat/merchStore/internal/config"
	v1 "github.com/magmaheat/merchStore/internal/http/v1"
	"github.com/magmaheat/merchStore/internal/repo/pgdb"
	"github.com/magmaheat/merchStore/internal/service"
	"github.com/magmaheat/merchStore/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

func Run() {
	cfg := config.New()

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

	_ = routes

}
