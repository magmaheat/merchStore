package app

import (
	"github.com/magmaheat/merchStore/internal/config"
	"github.com/magmaheat/merchStore/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

func Run() {
	cfg := config.New()

	setupLogger(cfg.Log.Level)

	log.Info("Initializing postgres")
	pg, err := postgres.New()
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pg.Close()

	log.Info("Initializing repositories...")

}
